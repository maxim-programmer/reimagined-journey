<template>
  <div class="page">
    <header class="page__header">
      <div class="logo">
        <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
          <rect width="28" height="28" rx="8" fill="#4f6ef7"/>
          <path d="M8 20V10l6-4 6 4v10" stroke="#fff" stroke-width="1.8" stroke-linejoin="round"/>
          <rect x="11" y="14" width="6" height="6" rx="1" stroke="#fff" stroke-width="1.6"/>
        </svg>
        <span class="logo__text">KnowledgeBase</span>
      </div>
    </header>

    <main class="page__main">
      <div class="hero">
        <h1 class="hero__title">База знаний</h1>
        <p class="hero__desc">Загружайте документы PDF и DOCX — система извлечёт текст и проиндексирует его для мгновенного полнотекстового поиска.</p>
      </div>

      <section class="upload-section">
        <DropZone @files-selected="onFilesSelected" />
        <FileQueue
          :files="fileQueue"
          :is-uploading="isUploading"
          @remove="removeFile"
          @upload-all="uploadAll"
          @clear="clearQueue"
        />
      </section>

      <div v-if="uploadedDocs.length > 0" class="docs-section">
        <div class="docs-section__header">
          <h2 class="docs-section__title">Загруженные документы</h2>
          <span class="docs-section__count">{{ uploadedDocs.length }}</span>
        </div>
        <ul class="docs-list">
          <li v-for="doc in uploadedDocs" :key="doc.id" class="doc-card">
            <div class="doc-card__icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
              </svg>
            </div>
            <div class="doc-card__info">
              <span class="doc-card__name">{{ doc.file_name }}</span>
              <span class="doc-card__meta">{{ formatSize(doc.file_size) }} · {{ formatDate(doc.uploaded_at) }}</span>
            </div>
            <span class="doc-card__badge">{{ doc.mime_type === 'application/pdf' ? 'PDF' : 'DOCX' }}</span>
          </li>
        </ul>
      </div>
    </main>
  </div>
</template>

<script>
import DropZone from '../components/DropZone.vue'
import FileQueue from '../components/FileQueue.vue'
import { uploadDocument, listDocuments } from '../api/documents.js'

let idCounter = 0

export default {
  name: 'HomePage',
  components: { DropZone, FileQueue },

  data() {
    return {
      fileQueue: [],
      uploadedDocs: [],
      isUploading: false,
    }
  },

  async created() {
    await this.fetchDocs()
  },

  methods: {
    async fetchDocs() {
      try {
        this.uploadedDocs = await listDocuments()
      } catch {
        this.uploadedDocs = []
      }
    },

    onFilesSelected(files) {
      for (const file of files) {
        const alreadyQueued = this.fileQueue.some(
          (f) => f.file.name === file.name && f.file.size === file.size,
        )
        if (!alreadyQueued) {
          this.fileQueue.push({
            id: ++idCounter,
            file,
            status: 'pending',
            error: '',
          })
        }
      }
    },

    removeFile(id) {
      this.fileQueue = this.fileQueue.filter((f) => f.id !== id)
    },

    clearQueue() {
      this.fileQueue = this.fileQueue.filter((f) => f.status === 'uploading' || f.status === 'indexing')
    },

    async uploadAll() {
      const pending = this.fileQueue.filter((f) => f.status === 'pending')
      if (pending.length === 0) return

      this.isUploading = true

      await Promise.allSettled(
        pending.map((item) => this.uploadOne(item)),
      )

      this.isUploading = false
      await this.fetchDocs()
    },

    async uploadOne(item) {
      item.status = 'uploading'
      try {
        await uploadDocument(item.file)
        item.status = 'indexing'
        await this.simulateIndexing()
        item.status = 'done'
      } catch (err) {
        item.status = 'error'
        item.error = err.message
      }
    },

    simulateIndexing() {
      return new Promise((resolve) => setTimeout(resolve, 800))
    },

    formatSize(bytes) {
      if (bytes < 1024) return bytes + ' Б'
      if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' КБ'
      return (bytes / (1024 * 1024)).toFixed(1) + ' МБ'
    },

    formatDate(iso) {
      return new Date(iso).toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
      })
    },
  },
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page__header {
  padding: 16px 40px;
  border-bottom: 1px solid #1e2236;
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo__text {
  font-size: 16px;
  font-weight: 700;
  color: #e8eaf0;
  letter-spacing: -0.02em;
}

.page__main {
  flex: 1;
  max-width: 720px;
  margin: 0 auto;
  padding: 48px 24px 80px;
  width: 100%;
}

.hero {
  margin-bottom: 36px;
}

.hero__title {
  font-size: 32px;
  font-weight: 700;
  color: #e8eaf0;
  letter-spacing: -0.03em;
  margin-bottom: 10px;
}

.hero__desc {
  font-size: 15px;
  color: #7b82a0;
  line-height: 1.6;
  max-width: 540px;
}

.upload-section {
  margin-bottom: 48px;
}

.docs-section__header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.docs-section__title {
  font-size: 16px;
  font-weight: 600;
  color: #e8eaf0;
}

.docs-section__count {
  font-size: 12px;
  background: #232840;
  color: #7b82a0;
  padding: 2px 8px;
  border-radius: 20px;
}

.docs-list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.doc-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 10px;
  transition: border-color 0.15s;
}

.doc-card:hover {
  border-color: #2e3347;
}

.doc-card__icon {
  color: #4f6ef7;
  flex-shrink: 0;
}

.doc-card__info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.doc-card__name {
  font-size: 14px;
  color: #e8eaf0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.doc-card__meta {
  font-size: 12px;
  color: #4a5070;
}

.doc-card__badge {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.05em;
  color: #4f6ef7;
  background: rgba(79, 110, 247, 0.12);
  padding: 3px 8px;
  border-radius: 6px;
  flex-shrink: 0;
}
</style>