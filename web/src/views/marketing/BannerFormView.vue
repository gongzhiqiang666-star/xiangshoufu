<template>
  <div class="banner-form">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑滚动图' : '新增滚动图' }}</span>
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

        <el-form-item label="图片" prop="image_url">
          <div class="image-uploader">
            <el-upload
              class="avatar-uploader"
              :show-file-list="false"
              :http-request="handleUpload"
              accept="image/jpeg,image/png,image/webp"
            >
              <img v-if="formData.image_url" :src="formData.image_url" class="uploaded-image" />
              <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
            </el-upload>
            <div class="upload-tip">建议尺寸: 750x400px，最大2MB</div>
          </div>
        </el-form-item>

        <el-form-item label="链接类型" prop="link_type">
          <el-radio-group v-model="formData.link_type">
            <el-radio v-for="opt in linkTypeOptions" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="formData.link_type !== LinkType.None" label="跳转链接" prop="link_url">
          <el-input
            v-model="formData.link_url"
            :placeholder="formData.link_type === LinkType.Internal ? '请输入内部页面路径' : '请输入外部链接URL'"
            maxlength="500"
          />
        </el-form-item>

        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" />
          <span class="form-tip">数值越大越靠前</span>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="展示时间">
          <el-date-picker
            v-model="timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
          <div class="form-tip">不设置则为长期有效</div>
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules, UploadRequestOptions } from 'element-plus'
import { getBannerDetail, createBanner, updateBanner } from '@/api/banner'
import { uploadImage } from '@/api/upload'
import type { BannerCreateRequest } from '@/types/banner'
import { LinkType, linkTypeOptions } from '@/types/banner'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const bannerId = computed(() => Number(route.params.id))

const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const timeRange = ref<[string, string] | null>(null)

const formData = reactive<BannerCreateRequest>({
  title: '',
  image_url: '',
  link_type: LinkType.None,
  link_url: '',
  sort_order: 0,
  status: 1,
  start_time: null,
  end_time: null
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { max: 100, message: '标题最多100个字符', trigger: 'blur' }
  ],
  image_url: [
    { required: true, message: '请上传图片', trigger: 'change' }
  ],
  link_url: [
    { max: 500, message: '链接最多500个字符', trigger: 'blur' }
  ]
}

// 监听时间范围变化
watch(timeRange, (val) => {
  if (val && val.length === 2) {
    formData.start_time = val[0]
    formData.end_time = val[1]
  } else {
    formData.start_time = null
    formData.end_time = null
  }
})

// 上传图片
const handleUpload = async (options: UploadRequestOptions) => {
  try {
    const res = await uploadImage(options.file as File, 'banner')
    formData.image_url = res.file_url
    ElMessage.success('上传成功')
  } catch (error) {
    ElMessage.error('上传失败')
  }
}

// 获取详情
const fetchDetail = async () => {
  if (!bannerId.value) return

  try {
    const res = await getBannerDetail(bannerId.value)
    Object.assign(formData, {
      title: res.title,
      image_url: res.image_url,
      link_type: res.link_type,
      link_url: res.link_url,
      sort_order: res.sort_order,
      status: res.status,
      start_time: res.start_time,
      end_time: res.end_time
    })

    // 设置时间范围
    if (res.start_time && res.end_time) {
      timeRange.value = [res.start_time, res.end_time]
    }
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
      await updateBanner(bannerId.value, formData)
      ElMessage.success('保存成功')
    } else {
      await createBanner(formData)
      ElMessage.success('创建成功')
    }
    router.push('/marketing/banners')
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
  router.push('/marketing/banners')
}

onMounted(() => {
  if (isEdit.value) {
    fetchDetail()
  }
})
</script>

<style scoped lang="scss">
.banner-form {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.image-uploader {
  .avatar-uploader {
    :deep(.el-upload) {
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      cursor: pointer;
      position: relative;
      overflow: hidden;
      transition: border-color 0.3s;
      width: 300px;
      height: 150px;
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

  .avatar-uploader-icon {
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
