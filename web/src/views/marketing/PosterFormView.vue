<template>
  <div class="poster-form">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑海报' : '新增海报' }}</span>
          <el-button @click="handleBack">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
        style="max-width: 600px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入标题" maxlength="100" show-word-limit />
        </el-form-item>

        <el-form-item label="分类" prop="category_id">
          <el-select v-model="formData.category_id" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="海报图片" prop="image_url">
          <div class="image-uploader">
            <el-upload
              class="poster-uploader"
              :show-file-list="false"
              :http-request="handleUpload"
              accept="image/jpeg,image/png,image/webp"
            >
              <img v-if="formData.image_url" :src="formData.thumbnail_url || formData.image_url" class="uploaded-image" />
              <el-icon v-else class="uploader-icon"><Plus /></el-icon>
            </el-upload>
            <div class="upload-tip">建议尺寸: 竖版海报，最大5MB</div>
          </div>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述（选填）"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" />
          <span class="form-tip">数值越大越靠前</span>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
            {{ isEdit ? '保存' : '创建' }}
          </el-button>
          <el-button @click="handleBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules, UploadRequestOptions } from 'element-plus'
import { getPosterDetail, createPoster, updatePoster, getPosterCategories } from '@/api/poster'
import { uploadImage } from '@/api/upload'
import type { PosterCreateRequest, PosterCategory } from '@/types/poster'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const posterId = computed(() => Number(route.params.id))

const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const categories = ref<PosterCategory[]>([])

const formData = reactive<PosterCreateRequest>({
  title: '',
  category_id: 0,
  image_url: '',
  thumbnail_url: '',
  description: '',
  file_size: 0,
  width: 0,
  height: 0,
  sort_order: 0,
  status: 1
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { max: 100, message: '标题最多100个字符', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  image_url: [
    { required: true, message: '请上传海报图片', trigger: 'change' }
  ]
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const res = await getPosterCategories({ status: 1 })
    categories.value = res || []
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 上传图片
const handleUpload = async (options: UploadRequestOptions) => {
  try {
    const res = await uploadImage(options.file as File, 'poster')
    formData.image_url = res.file_url
    formData.thumbnail_url = res.thumbnail_url || ''
    formData.file_size = res.file_size
    formData.width = res.width
    formData.height = res.height
    ElMessage.success('上传成功')
  } catch (error) {
    ElMessage.error('上传失败')
  }
}

// 获取详情
const fetchDetail = async () => {
  if (!posterId.value) return

  try {
    const res = await getPosterDetail(posterId.value)
    Object.assign(formData, {
      title: res.title,
      category_id: res.category_id,
      image_url: res.image_url,
      thumbnail_url: res.thumbnail_url,
      description: res.description,
      file_size: res.file_size,
      width: res.width,
      height: res.height,
      sort_order: res.sort_order,
      status: res.status
    })
  } catch (error) {
    ElMessage.error('获取详情失败')
    router.back()
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitLoading.value = true

    if (isEdit.value) {
      await updatePoster(posterId.value, formData)
      ElMessage.success('保存成功')
    } else {
      await createPoster(formData)
      ElMessage.success('创建成功')
    }
    router.push('/marketing/posters')
  } catch (error: any) {
    if (error !== false) {
      ElMessage.error(isEdit.value ? '保存失败' : '创建失败')
    }
  } finally {
    submitLoading.value = false
  }
}

// 返回
const handleBack = () => {
  router.push('/marketing/posters')
}

onMounted(() => {
  fetchCategories()
  if (isEdit.value) {
    fetchDetail()
  }
})
</script>

<style scoped lang="scss">
.poster-form {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.image-uploader {
  .poster-uploader {
    :deep(.el-upload) {
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      cursor: pointer;
      position: relative;
      overflow: hidden;
      transition: border-color 0.3s;
      width: 200px;
      height: 300px;
      display: flex;
      align-items: center;
      justify-content: center;

      &:hover {
        border-color: #409eff;
      }
    }
  }

  .uploaded-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .uploader-icon {
    font-size: 28px;
    color: #8c939d;
  }

  .upload-tip {
    color: #909399;
    font-size: 12px;
    margin-top: 8px;
  }
}

.form-tip {
  color: #909399;
  font-size: 12px;
  margin-left: 10px;
}
</style>
