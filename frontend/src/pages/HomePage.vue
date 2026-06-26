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

        <div class="docs-table-wrap">
          <table class="docs-table">
            <thead>
              <tr>
                <th class="docs-table__th">Название</th>
                <th class="docs-table__th docs-table__th--center">Тип</th>
                <th class="docs-table__th">Дата загрузки</th>
                <th class="docs-table__th docs-table__th--center">Статус</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="doc in uploadedDocs"
                :key="doc.id"
                class="docs-table__row"
              >
                <td class="docs-table__td">
                  <div class="doc-name-cell">
                    <span class="doc-name-cell__icon">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                        <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                      </svg>
                    </span>
                    <span class="doc-name-cell__name" :title="doc.file_name">{{ doc.file_name }}</span>
                  </div>
                </td>
                <td class="docs-table__td docs-table__td--center">
                  <span class="badge badge--type">
                    {{ doc.mime_type === 'application/pdf' ? 'PDF' : 'DOCX' }}
                  </span>
                </td>
                <td class="docs-table__td docs-table__td--date">
                  {{ formatDate(doc.uploaded_at) }}
                </td>
                <td class="docs-table__td docs-table__td--center">
                  <span class="badge" :class="statusBadgeClass(doc.status)">
                    {{ statusLabel(doc.status) }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
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

    statusLabel(status) {
      const map = {
        uploaded: 'Загружен',
        processing: 'Обработка',
        indexed: 'Индексирован',
        error: 'Ошибка',
      }
      return map[status] ?? status
    },

    statusBadgeClass(status) {
      const map = {
        uploaded: 'badge--uploaded',
        processing: 'badge--processing',
        indexed: 'badge--indexed',
        error: 'badge--error',
      }
      return map[status] ?? 'badge--uploaded'
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

.docs-table-wrap {
  border: 1px solid #232840;
  border-radius: 12px;
  overflow: hidden;
}

.docs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.docs-table__th {
  padding: 11px 16px;
  text-align: left;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #7b82a0;
  background: #161a27;
  border-bottom: 1px solid #232840;
}

.docs-table__th--center {
  text-align: center;
}

.docs-table__row {
  background: #161a27;
  transition: background 0.15s;
}

.docs-table__row:not(:last-child) {
  border-bottom: 1px solid #1e2236;
}

.docs-table__row:hover {
  background: #1b2035;
}

.docs-table__td {
  padding: 12px 16px;
  color: #e8eaf0;
  vertical-align: middle;
}

.docs-table__td--center {
  text-align: center;
}

.docs-table__td--date {
  color: #7b82a0;
  white-space: nowrap;
  font-size: 13px;
}

.doc-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.doc-name-cell__icon {
  color: #4f6ef7;
  flex-shrink: 0;
  display: flex;
}

.doc-name-cell__name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 280px;
}

.badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.04em;
  padding: 3px 9px;
  border-radius: 6px;
  white-space: nowrap;
}

.badge--type {
  color: #4f6ef7;
  background: rgba(79, 110, 247, 0.12);
}

.badge--uploaded {
  color: #7b82a0;
  background: rgba(123, 130, 160, 0.12);
}

.badge--processing {
  color: #f7a24f;
  background: rgba(247, 162, 79, 0.12);
}

.badge--indexed {
  color: #3ec97a;
  background: rgba(62, 201, 122, 0.12);
}

.badge--error {
  color: #e05c5c;
  background: rgba(224, 92, 92, 0.12);
}
</style>