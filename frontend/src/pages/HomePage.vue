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

      <section class="search-section">
        <div class="search-bar">
          <div class="search-bar__input-wrap">
            <svg class="search-bar__icon" width="18" height="18" viewBox="0 0 24 24" fill="none">
              <circle cx="11" cy="11" r="7" stroke="currentColor" stroke-width="1.8"/>
              <line x1="16.5" y1="16.5" x2="21" y2="21" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
            </svg>
            <input
              v-model="searchQuery"
              class="search-bar__input"
              type="text"
              placeholder="Введите поисковый запрос…"
              @keydown.enter="onSearch"
            />
            <button
              v-if="searchQuery"
              class="search-bar__clear"
              @click="clearSearch"
              title="Очистить"
            >
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                <line x1="2" y1="2" x2="12" y2="12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                <line x1="12" y1="2" x2="2" y2="12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
              </svg>
            </button>
          </div>
          <button
            class="search-bar__btn"
            :disabled="isSearching || !searchQuery.trim()"
            @click="onSearch"
          >
            <svg v-if="isSearching" class="spin" width="16" height="16" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" stroke-dasharray="28 56" stroke-linecap="round"/>
            </svg>
            <span>{{ isSearching ? 'Поиск…' : 'Найти' }}</span>
          </button>
        </div>

        <div v-if="searchError" class="search-error">
          {{ searchError }}
        </div>

        <div v-if="searchDone" class="search-results">
          <div class="search-results__header">
            <span class="search-results__label">
              {{ searchHits.length > 0 ? `Найдено фрагментов: ${searchHits.length}` : 'Ничего не найдено' }}
            </span>
            <span v-if="lastQuery" class="search-results__query">«{{ lastQuery }}»</span>
          </div>

          <ul v-if="searchHits.length > 0" class="search-hits">
            <li
              v-for="(hit, index) in searchHits"
              :key="hit.chunk_id || index"
              class="search-hit"
            >
              <div class="search-hit__meta">
                <span class="search-hit__filename">{{ hit.file_name }}</span>
                <span v-if="hit.page" class="search-hit__page">стр. {{ hit.page }}</span>
                <span class="search-hit__score">релевантность: {{ hit.score.toFixed(2) }}</span>
              </div>
              <p class="search-hit__text">{{ hit.text }}</p>
            </li>
          </ul>

          <div v-else class="search-empty">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none">
              <circle cx="11" cy="11" r="7" stroke="currentColor" stroke-width="1.5"/>
              <line x1="16.5" y1="16.5" x2="21" y2="21" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
            <p>По вашему запросу ничего не найдено. Попробуйте изменить запрос.</p>
          </div>
        </div>
      </section>

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
import { uploadDocument, listDocuments, searchDocuments } from '../api/documents.js'

let idCounter = 0

export default {
  name: 'HomePage',
  components: { DropZone, FileQueue },

  data() {
    return {
      fileQueue: [],
      uploadedDocs: [],
      isUploading: false,
      searchQuery: '',
      isSearching: false,
      searchHits: [],
      searchDone: false,
      searchError: '',
      lastQuery: '',
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

    async onSearch() {
      const query = this.searchQuery.trim()
      if (!query || this.isSearching) return

      this.isSearching = true
      this.searchError = ''
      this.searchDone = false
      this.searchHits = []

      try {
        this.searchHits = await searchDocuments(query)
        this.lastQuery = query
        this.searchDone = true
      } catch (err) {
        this.searchError = err.message
      } finally {
        this.isSearching = false
      }
    },

    clearSearch() {
      this.searchQuery = ''
      this.searchHits = []
      this.searchDone = false
      this.searchError = ''
      this.lastQuery = ''
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

.search-section {
  margin-bottom: 48px;
}

.search-bar {
  display: flex;
  gap: 10px;
}

.search-bar__input-wrap {
  position: relative;
  flex: 1;
  display: flex;
  align-items: center;
}

.search-bar__icon {
  position: absolute;
  left: 14px;
  color: #4a5070;
  pointer-events: none;
  flex-shrink: 0;
}

.search-bar__input {
  width: 100%;
  background: #161a27;
  border: 1.5px solid #2e3347;
  border-radius: 10px;
  padding: 11px 40px 11px 42px;
  font-size: 14px;
  color: #e8eaf0;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}

.search-bar__input::placeholder {
  color: #4a5070;
}

.search-bar__input:focus {
  border-color: #4f6ef7;
}

.search-bar__clear {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  cursor: pointer;
  color: #4a5070;
  padding: 4px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.15s;
}

.search-bar__clear:hover {
  color: #e8eaf0;
}

.search-bar__btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 22px;
  background: #4f6ef7;
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.15s, opacity 0.15s;
}

.search-bar__btn:not(:disabled):hover {
  background: #3d5ce8;
}

.search-bar__btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.search-error {
  margin-top: 12px;
  font-size: 13px;
  color: #e05c5c;
  padding: 10px 14px;
  background: rgba(224, 92, 92, 0.08);
  border: 1px solid rgba(224, 92, 92, 0.2);
  border-radius: 8px;
}

.search-results {
  margin-top: 16px;
}

.search-results__header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.search-results__label {
  font-size: 13px;
  font-weight: 600;
  color: #7b82a0;
}

.search-results__query {
  font-size: 13px;
  color: #4a5070;
}

.search-hits {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.search-hit {
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 10px;
  padding: 14px 16px;
  transition: border-color 0.15s;
}

.search-hit:hover {
  border-color: #4f6ef7;
}

.search-hit__meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.search-hit__filename {
  font-size: 13px;
  font-weight: 600;
  color: #4f6ef7;
}

.search-hit__page {
  font-size: 12px;
  color: #7b82a0;
  background: #1e2236;
  padding: 2px 8px;
  border-radius: 4px;
}

.search-hit__score {
  font-size: 11px;
  color: #4a5070;
  margin-left: auto;
}

.search-hit__text {
  font-size: 13px;
  color: #c0c4d6;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.search-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 36px 0;
  color: #4a5070;
  text-align: center;
}

.search-empty p {
  font-size: 14px;
  color: #4a5070;
  max-width: 340px;
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

.spin {
  animation: spin 0.9s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>