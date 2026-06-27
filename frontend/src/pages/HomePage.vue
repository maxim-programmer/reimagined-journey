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
      <div class="header-right">
        <span class="header-user">{{ user.login }}</span>
        <button class="btn-logout" @click="onLogout">Выйти</button>
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
              @focus="showHistory = history.length > 0"
              @blur="onSearchBlur"
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

            <div v-if="showHistory && history.length > 0" class="search-history-dropdown">
              <div class="search-history-dropdown__header">
                <span class="search-history-dropdown__label">История поиска</span>
                <button class="search-history-dropdown__clear" @mousedown.prevent="onClearHistory">
                  Очистить
                </button>
              </div>
              <ul class="search-history-dropdown__list">
                <li
                  v-for="item in history"
                  :key="item.id"
                  class="search-history-dropdown__item"
                  @mousedown.prevent="pickHistory(item.query)"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="1.8"/>
                    <polyline points="12 7 12 12 15 14" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                  </svg>
                  <span>{{ item.query }}</span>
                </li>
              </ul>
            </div>
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

          <ul v-if="pagedHits.length > 0" class="search-hits">
            <li
              v-for="(hit, index) in pagedHits"
              :key="hit.chunk_id || index"
              class="hit-card"
            >
              <div class="hit-card__header">
                <div class="hit-card__file">
                  <svg class="hit-card__file-icon" width="14" height="14" viewBox="0 0 24 24" fill="none">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                    <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                  </svg>
                  <span class="hit-card__filename">{{ hit.file_name }}</span>
                </div>
                <div class="hit-card__score">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
                    <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                  </svg>
                  <span>{{ hit.score.toFixed(2) }}</span>
                </div>
              </div>

              <div class="hit-card__divider"></div>

              <div class="hit-card__meta">
                <div v-if="hit.page" class="hit-card__page">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                  </svg>
                  <span>Страница {{ hit.page }}</span>
                </div>
                <div v-else class="hit-card__page hit-card__page--unknown">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                  </svg>
                  <span>Страница не определена</span>
                </div>
              </div>

              <p class="hit-card__text" v-html="highlight(hit.text)"></p>

              <div class="hit-card__footer">
                <span class="hit-card__relevance-label">Релевантность</span>
                <div class="hit-card__relevance-bar-wrap">
                  <div
                    class="hit-card__relevance-bar"
                    :style="{ width: relevancePercent(hit.score) }"
                  ></div>
                </div>
              </div>
            </li>
          </ul>

          <div v-if="searchHits.length > 0" class="pagination">
            <button
              class="pagination__btn"
              :disabled="currentPage === 1"
              @click="goToPage(currentPage - 1)"
              title="Предыдущая страница"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <polyline points="15 18 9 12 15 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </button>
            <div class="pagination__pages">
              <button
                v-for="page in visiblePages"
                :key="page"
                class="pagination__page"
                :class="{
                  'pagination__page--active': page === currentPage,
                  'pagination__page--ellipsis': page === '…',
                }"
                :disabled="page === '…'"
                @click="page !== '…' && goToPage(page)"
              >
                {{ page }}
              </button>
            </div>
            <button
              class="pagination__btn"
              :disabled="currentPage === totalPages"
              @click="goToPage(currentPage + 1)"
              title="Следующая страница"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <polyline points="9 18 15 12 9 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </button>
          </div>

          <div v-if="searchHits.length === 0" class="search-empty">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none">
              <circle cx="11" cy="11" r="7" stroke="currentColor" stroke-width="1.5"/>
              <line x1="16.5" y1="16.5" x2="21" y2="21" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
            <p>По вашему запросу ничего не найдено. Попробуйте изменить формулировку.</p>
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
                <th class="docs-table__th docs-table__th--hide-sm">Дата загрузки</th>
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
                <td class="docs-table__td docs-table__td--date docs-table__td--hide-sm">
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
import { uploadDocument, listDocuments, searchDocuments, getSearchHistory, clearSearchHistory, logout } from '../api/documents.js'

const PAGE_SIZE = 10
let idCounter = 0

export default {
  name: 'HomePage',
  components: { DropZone, FileQueue },
  emits: ['logout'],

  props: {
    user: {
      type: Object,
      required: true,
    },
  },

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
      currentPage: 1,
      history: [],
      showHistory: false,
    }
  },

  computed: {
    totalPages() {
      return Math.ceil(this.searchHits.length / PAGE_SIZE)
    },

    pagedHits() {
      const start = (this.currentPage - 1) * PAGE_SIZE
      return this.searchHits.slice(start, start + PAGE_SIZE)
    },

    visiblePages() {
      const total = this.totalPages
      const current = this.currentPage
      if (total <= 7) {
        return Array.from({ length: total }, (_, i) => i + 1)
      }
      const pages = []
      if (current <= 4) {
        for (let i = 1; i <= 5; i++) pages.push(i)
        pages.push('…')
        pages.push(total)
      } else if (current >= total - 3) {
        pages.push(1)
        pages.push('…')
        for (let i = total - 4; i <= total; i++) pages.push(i)
      } else {
        pages.push(1)
        pages.push('…')
        pages.push(current - 1)
        pages.push(current)
        pages.push(current + 1)
        pages.push('…')
        pages.push(total)
      }
      return pages
    },
  },

  async created() {
    await Promise.all([this.fetchDocs(), this.fetchHistory()])
  },

  methods: {
    async fetchDocs() {
      try {
        this.uploadedDocs = await listDocuments()
      } catch {
        this.uploadedDocs = []
      }
    },

    async fetchHistory() {
      try {
        this.history = await getSearchHistory()
      } catch {
        this.history = []
      }
    },

    async onSearch() {
      const query = this.searchQuery.trim()
      if (!query || this.isSearching) return

      this.showHistory = false
      this.isSearching = true
      this.searchError = ''
      this.searchDone = false
      this.searchHits = []
      this.currentPage = 1

      try {
        this.searchHits = await searchDocuments(query)
        this.lastQuery = query
        this.searchDone = true
        await this.fetchHistory()
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
      this.currentPage = 1
      this.showHistory = false
    },

    onSearchBlur() {
      setTimeout(() => { this.showHistory = false }, 150)
    },

    pickHistory(query) {
      this.searchQuery = query
      this.showHistory = false
      this.onSearch()
    },

    async onClearHistory() {
      try {
        await clearSearchHistory()
        this.history = []
        this.showHistory = false
      } catch {}
    },

    goToPage(page) {
      if (page < 1 || page > this.totalPages) return
      this.currentPage = page
      const el = document.querySelector('.search-results')
      if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    },

    highlight(text) {
      if (!this.lastQuery || !text) return this.escapeHtml(text)
      const words = this.lastQuery.trim().split(/\s+/).filter(Boolean).map((w) => this.escapeRegex(w))
      if (words.length === 0) return this.escapeHtml(text)
      const pattern = new RegExp(`(${words.join('|')})`, 'gi')
      return this.escapeHtml(text).replace(pattern, '<mark class="search-highlight">$1</mark>')
    },

    escapeHtml(str) {
      return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;')
    },

    escapeRegex(str) {
      return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    },

    relevancePercent(score) {
      const maxScore = Math.max(...this.searchHits.map((h) => h.score), 1)
      const pct = Math.min((score / maxScore) * 100, 100)
      return pct.toFixed(1) + '%'
    },

    onFilesSelected(files) {
      for (const file of files) {
        const alreadyQueued = this.fileQueue.some(
          (f) => f.file.name === file.name && f.file.size === file.size,
        )
        if (!alreadyQueued) {
          this.fileQueue.push({ id: ++idCounter, file, status: 'pending', error: '' })
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
      await Promise.allSettled(pending.map((item) => this.uploadOne(item)))
      this.isUploading = false
      await this.fetchDocs()
    },

    async uploadOne(item) {
      item.status = 'uploading'
      try {
        await uploadDocument(item.file)
        item.status = 'indexing'
        await new Promise((resolve) => setTimeout(resolve, 800))
        item.status = 'done'
      } catch (err) {
        item.status = 'error'
        item.error = err.message
      }
    },

    async onLogout() {
      await logout()
      this.$emit('logout')
    },

    statusLabel(status) {
      const map = { uploaded: 'Загружен', processing: 'Обработка', indexed: 'Индексирован', error: 'Ошибка' }
      return map[status] ?? status
    },

    statusBadgeClass(status) {
      const map = { uploaded: 'badge--uploaded', processing: 'badge--processing', indexed: 'badge--indexed', error: 'badge--error' }
      return map[status] ?? 'badge--uploaded'
    },

    formatDate(iso) {
      return new Date(iso).toLocaleDateString('ru-RU', { day: '2-digit', month: 'short', year: 'numeric' })
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
  justify-content: space-between;
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

.header-right {
  display: flex;
  align-items: center;
  gap: 14px;
}

.header-user {
  font-size: 13px;
  color: #7b82a0;
}

.btn-logout {
  background: none;
  border: 1px solid #2e3347;
  border-radius: 8px;
  color: #7b82a0;
  font-size: 13px;
  font-family: inherit;
  padding: 6px 14px;
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}

.btn-logout:hover {
  border-color: #e05c5c;
  color: #e05c5c;
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
  z-index: 1;
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
  z-index: 1;
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

.search-history-dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  background: #1e2236;
  border: 1px solid #2e3347;
  border-radius: 10px;
  z-index: 100;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.search-history-dropdown__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px 8px;
  border-bottom: 1px solid #2e3347;
}

.search-history-dropdown__label {
  font-size: 11px;
  font-weight: 600;
  color: #4a5070;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.search-history-dropdown__clear {
  background: none;
  border: none;
  font-size: 12px;
  color: #4a5070;
  cursor: pointer;
  font-family: inherit;
  padding: 2px 6px;
  border-radius: 4px;
  transition: color 0.15s;
}

.search-history-dropdown__clear:hover {
  color: #e05c5c;
}

.search-history-dropdown__list {
  list-style: none;
  max-height: 220px;
  overflow-y: auto;
}

.search-history-dropdown__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  font-size: 13px;
  color: #c0c4d6;
  cursor: pointer;
  transition: background 0.12s;
}

.search-history-dropdown__item:hover {
  background: #232840;
}

.search-history-dropdown__item svg {
  color: #4a5070;
  flex-shrink: 0;
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
  margin-bottom: 14px;
  flex-wrap: wrap;
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
  gap: 12px;
}

.hit-card {
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 12px;
  overflow: hidden;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.hit-card:hover {
  border-color: #4f6ef7;
  box-shadow: 0 0 0 3px rgba(79, 110, 247, 0.08);
}

.hit-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px 12px;
  gap: 12px;
}

.hit-card__file {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.hit-card__file-icon {
  color: #4f6ef7;
  flex-shrink: 0;
}

.hit-card__filename {
  font-size: 14px;
  font-weight: 600;
  color: #4f6ef7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.hit-card__score {
  display: flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
  font-size: 12px;
  font-weight: 600;
  color: #f7c14f;
  background: rgba(247, 193, 79, 0.1);
  border: 1px solid rgba(247, 193, 79, 0.2);
  border-radius: 6px;
  padding: 3px 9px;
}

.hit-card__divider {
  height: 1px;
  background: #1e2236;
  margin: 0 16px;
}

.hit-card__meta {
  padding: 10px 16px 0;
}

.hit-card__page {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #7b82a0;
  background: #1e2236;
  border-radius: 5px;
  padding: 3px 9px;
}

.hit-card__page--unknown {
  color: #4a5070;
}

.hit-card__text {
  font-size: 13px;
  color: #c0c4d6;
  line-height: 1.7;
  padding: 10px 16px 0;
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.hit-card__footer {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px 14px;
  margin-top: 10px;
}

.hit-card__relevance-label {
  font-size: 11px;
  color: #4a5070;
  white-space: nowrap;
  flex-shrink: 0;
}

.hit-card__relevance-bar-wrap {
  flex: 1;
  height: 3px;
  background: #1e2236;
  border-radius: 2px;
  overflow: hidden;
}

.hit-card__relevance-bar {
  height: 100%;
  background: linear-gradient(90deg, #4f6ef7, #7b96ff);
  border-radius: 2px;
  transition: width 0.4s ease;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 24px;
  flex-wrap: wrap;
}

.pagination__btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 8px;
  color: #7b82a0;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s, color 0.15s;
  flex-shrink: 0;
}

.pagination__btn:not(:disabled):hover {
  background: #1b2035;
  border-color: #4f6ef7;
  color: #e8eaf0;
}

.pagination__btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.pagination__pages {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
  justify-content: center;
}

.pagination__page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  padding: 0 6px;
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  color: #7b82a0;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s, color 0.15s;
}

.pagination__page:not(:disabled):not(.pagination__page--active):hover {
  background: #1b2035;
  border-color: #4f6ef7;
  color: #e8eaf0;
}

.pagination__page--active {
  background: #4f6ef7;
  border-color: #4f6ef7;
  color: #fff;
  cursor: default;
}

.pagination__page--ellipsis {
  border-color: transparent;
  background: transparent;
  cursor: default;
  color: #4a5070;
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
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.docs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
  min-width: 400px;
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
  white-space: nowrap;
}

.docs-table__th--center { text-align: center; }

.docs-table__row {
  background: #161a27;
  transition: background 0.15s;
}

.docs-table__row:not(:last-child) {
  border-bottom: 1px solid #1e2236;
}

.docs-table__row:hover { background: #1b2035; }

.docs-table__td {
  padding: 12px 16px;
  color: #e8eaf0;
  vertical-align: middle;
}

.docs-table__td--center { text-align: center; }

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

.badge--type { color: #4f6ef7; background: rgba(79, 110, 247, 0.12); }
.badge--uploaded { color: #7b82a0; background: rgba(123, 130, 160, 0.12); }
.badge--processing { color: #f7a24f; background: rgba(247, 162, 79, 0.12); }
.badge--indexed { color: #3ec97a; background: rgba(62, 201, 122, 0.12); }
.badge--error { color: #e05c5c; background: rgba(224, 92, 92, 0.12); }

.spin { animation: spin 0.9s linear infinite; }

@keyframes spin { to { transform: rotate(360deg); } }

@media (min-width: 1280px) {
  .page__main { max-width: 960px; }
  .hero__title { font-size: 40px; }
  .doc-name-cell__name { max-width: 480px; }
}

@media (min-width: 1920px) {
  .page__main { max-width: 1200px; }
}

@media (max-width: 768px) {
  .page__header { padding: 14px 16px; }
  .page__main { padding: 32px 16px 60px; }
  .hero { margin-bottom: 28px; }
  .hero__title { font-size: 24px; }
  .hero__desc { font-size: 14px; }
  .search-bar { flex-direction: column; gap: 8px; }
  .search-bar__btn { padding: 11px 20px; justify-content: center; }
  .docs-table__th--hide-sm, .docs-table__td--hide-sm { display: none; }
  .doc-name-cell__name { max-width: 160px; }
  .header-user { display: none; }
}

@media (max-width: 480px) {
  .hero__title { font-size: 20px; }
  .hit-card__header { flex-wrap: wrap; gap: 8px; }
  .pagination__page, .pagination__btn { min-width: 32px; width: 32px; height: 32px; font-size: 12px; }
}

@media (max-width: 360px) {
  .page__main { padding: 24px 12px 48px; }
  .search-bar__input { font-size: 13px; padding: 10px 36px 10px 38px; }
  .hit-card__text { font-size: 12px; }
  .doc-name-cell__name { max-width: 100px; }
}
</style>

<style>
.search-highlight {
  background-color: #f7c14f;
  color: #0f1117;
  border-radius: 2px;
  padding: 0 2px;
}
</style>