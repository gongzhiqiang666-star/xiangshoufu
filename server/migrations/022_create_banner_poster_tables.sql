-- 滚动图和营销海报功能数据库迁移
-- Migration: 022_create_banner_poster_tables.sql
-- Description: 创建滚动图、海报分类、营销海报、文件上传记录表

-- 滚动图表
CREATE TABLE IF NOT EXISTS banners (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(100) NOT NULL,           -- 标题
    image_url       VARCHAR(500) NOT NULL,           -- 图片URL
    link_type       SMALLINT DEFAULT 0,              -- 0无链接 1内部页面 2外部链接
    link_url        VARCHAR(500),                    -- 跳转链接
    sort_order      INT DEFAULT 0,                   -- 排序（越大越靠前）
    status          SMALLINT DEFAULT 1,              -- 1启用 0禁用
    start_time      TIMESTAMPTZ,                     -- 开始展示时间
    end_time        TIMESTAMPTZ,                     -- 结束展示时间
    click_count     BIGINT DEFAULT 0,                -- 点击统计
    created_by      BIGINT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 海报分类表
CREATE TABLE IF NOT EXISTS poster_categories (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(50) NOT NULL,            -- 分类名称
    sort_order      INT DEFAULT 0,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 营销海报表
CREATE TABLE IF NOT EXISTS posters (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(100) NOT NULL,           -- 标题
    category_id     BIGINT NOT NULL REFERENCES poster_categories(id),
    image_url       VARCHAR(500) NOT NULL,           -- 原图URL
    thumbnail_url   VARCHAR(500),                    -- 缩略图URL
    description     TEXT,                            -- 描述
    file_size       BIGINT DEFAULT 0,                -- 文件大小（字节）
    width           INT DEFAULT 0,                   -- 图片宽度
    height          INT DEFAULT 0,                   -- 图片高度
    sort_order      INT DEFAULT 0,
    status          SMALLINT DEFAULT 1,
    download_count  BIGINT DEFAULT 0,                -- 下载次数
    share_count     BIGINT DEFAULT 0,                -- 分享次数
    created_by      BIGINT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 文件上传记录表（通用）
CREATE TABLE IF NOT EXISTS uploaded_files (
    id              BIGSERIAL PRIMARY KEY,
    original_name   VARCHAR(255) NOT NULL,
    stored_name     VARCHAR(255) NOT NULL,
    file_path       VARCHAR(500) NOT NULL,
    file_url        VARCHAR(500) NOT NULL,
    file_size       BIGINT DEFAULT 0,
    mime_type       VARCHAR(100),
    width           INT,
    height          INT,
    module          VARCHAR(50),                     -- banner/poster
    ref_id          BIGINT,
    uploaded_by     BIGINT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_banners_status ON banners(status);
CREATE INDEX IF NOT EXISTS idx_banners_sort_order ON banners(sort_order DESC);
CREATE INDEX IF NOT EXISTS idx_banners_time_range ON banners(start_time, end_time);

CREATE INDEX IF NOT EXISTS idx_poster_categories_status ON poster_categories(status);
CREATE INDEX IF NOT EXISTS idx_poster_categories_sort_order ON poster_categories(sort_order DESC);

CREATE INDEX IF NOT EXISTS idx_posters_category_id ON posters(category_id);
CREATE INDEX IF NOT EXISTS idx_posters_status ON posters(status);
CREATE INDEX IF NOT EXISTS idx_posters_sort_order ON posters(sort_order DESC);

CREATE INDEX IF NOT EXISTS idx_uploaded_files_module ON uploaded_files(module);
CREATE INDEX IF NOT EXISTS idx_uploaded_files_ref_id ON uploaded_files(ref_id);

-- 注释
COMMENT ON TABLE banners IS '滚动图表';
COMMENT ON COLUMN banners.link_type IS '链接类型: 0无链接 1内部页面 2外部链接';
COMMENT ON COLUMN banners.status IS '状态: 1启用 0禁用';

COMMENT ON TABLE poster_categories IS '海报分类表';
COMMENT ON TABLE posters IS '营销海报表';
COMMENT ON TABLE uploaded_files IS '文件上传记录表';
