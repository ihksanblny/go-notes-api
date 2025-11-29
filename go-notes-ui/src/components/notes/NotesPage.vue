<script setup>
import { ref, onMounted, watch } from 'vue'
import { getNotes, createNote, deleteNote, updateNote } from '../../api/notes'
import NoteForm from './NoteForm.vue'
import NotesList from './NotesList.vue'

const notes = ref([])
const loading = ref(false)
const error = ref('')
const editingNote = ref(null)

// Pagination & Search State
const page = ref(1)
const limit = ref(6) // 6 items per page for grid layout
const total = ref(0)
const searchQuery = ref('')
const totalPages = ref(1)

// Debounce search
let searchTimeout = null

async function loadNotes() {
  loading.value = true
  error.value = ''
  try {
    const res = await getNotes({
      page: page.value,
      limit: limit.value,
      q: searchQuery.value
    })
    
    // Handle response structure { data: [], total: N, page: N, limit: N }
    notes.value = res.data || []
    total.value = res.total || 0
    totalPages.value = Math.ceil(total.value / limit.value) || 1
    
  } catch (err) {
    error.value = err.message || 'Terjadi kesalahan'
    notes.value = []
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1 // Reset to first page on new search
    loadNotes()
  }, 300)
}

function changePage(newPage) {
  if (newPage >= 1 && newPage <= totalPages.value) {
    page.value = newPage
    loadNotes()
    // Scroll to top of list
    document.getElementById('notes-grid')?.scrollIntoView({ behavior: 'smooth' })
  }
}

async function handleCreateNote({ title, content }) {
  if (!title.trim()) {
    error.value = 'Title wajib diisi'
    return
  }

  try {
    error.value = ''
    await createNote({ title, content })
    // Reload to show new note (usually on first page)
    searchQuery.value = ''
    page.value = 1
    await loadNotes()
  } catch (err) {
    error.value = err.message || 'Gagal menambah note'
  }
}

async function handleDeleteNote(id) {
  const ok = window.confirm('Yakin mau hapus note ini?')
  if (!ok) return

  try {
    error.value = ''
    await deleteNote(id)
    await loadNotes() // Reload to update list & pagination
    if (editingNote.value && editingNote.value.id === id) {
      editingNote.value = null
    }
  } catch (err) {
    error.value = err.message || 'Gagal menghapus note'
  }
}

function handleEditNote(note) {
  editingNote.value = note
  error.value = ''
  // Scroll to form
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function handleCancelEdit() {
  editingNote.value = null
  error.value = ''
}

async function handleUpdateNote({ id, title, content }) {
  if (!title.trim()) {
    error.value = 'Title wajib diisi'
    return
  }

  try {
    error.value = ''
    await updateNote(id, { title, content })
    await loadNotes() // Reload to show updated data
    editingNote.value = null
  } catch (err) {
    error.value = err.message || 'Gagal mengupdate note'
  }
}

onMounted(loadNotes)
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] font-sans text-slate-900 selection:bg-primary-100 selection:text-primary-900">
    
    <!-- Navbar / Header -->
    <header class="sticky top-0 z-30 bg-white/80 backdrop-blur-md border-b border-slate-200/60">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div class="bg-primary-600 text-white p-1.5 rounded-lg">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h1 class="font-display font-bold text-xl tracking-tight text-slate-900">
            SuperNotes
          </h1>
        </div>

        <div class="flex items-center gap-4">
          <div class="hidden sm:flex items-center gap-2 text-xs font-medium text-slate-500 bg-slate-100 px-3 py-1.5 rounded-full">
            <span class="w-2 h-2 rounded-full" :class="loading ? 'bg-amber-400 animate-pulse' : 'bg-emerald-400'"></span>
            {{ loading ? 'Syncing...' : 'All saved' }}
          </div>
        </div>
      </div>
    </header>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
      
      <!-- Error Toast -->
      <transition
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="transform -translate-y-2 opacity-0"
        enter-to-class="transform translate-y-0 opacity-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="transform translate-y-0 opacity-100"
        leave-to-class="transform -translate-y-2 opacity-0"
      >
        <div v-if="error" class="mb-6 p-4 bg-rose-50 border border-rose-100 rounded-xl flex items-center gap-3 text-rose-700 shadow-sm">
          <svg class="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-sm font-medium">{{ error }}</span>
        </div>
      </transition>

      <div class="grid lg:grid-cols-[380px_1fr] gap-8 lg:gap-12 items-start">
        
        <!-- Left Column: Create Note -->
        <aside class="lg:sticky lg:top-24 space-y-8">
          <div class="space-y-2">
            <h2 class="font-display font-bold text-3xl text-slate-900 leading-tight">
              Capture your <br/>
              <span class="text-primary-600">best ideas.</span>
            </h2>
            <p class="text-slate-500 text-lg">
              Simple, fast, and distraction-free note taking.
            </p>
          </div>

          <NoteForm
            :note-to-edit="editingNote"
            @create-note="handleCreateNote"
            @update-note="handleUpdateNote"
            @cancel-edit="handleCancelEdit"
          />
        </aside>

        <!-- Right Column: Notes List -->
        <main class="space-y-6">
          <!-- Search & Filter Bar -->
          <div class="bg-white p-2 rounded-2xl shadow-sm border border-slate-200/60 flex items-center gap-2">
            <div class="pl-3 text-slate-400">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              v-model="searchQuery"
              @input="handleSearch"
              type="text"
              placeholder="Search your notes..."
              class="flex-1 bg-transparent border-none outline-none text-sm text-slate-700 placeholder:text-slate-400 focus:ring-0 py-2.5"
            />
            <div class="pr-2">
              <span class="text-xs font-bold text-slate-400 bg-slate-100 px-2 py-1 rounded-md">
                {{ total }} NOTES
              </span>
            </div>
          </div>

          <!-- Notes Grid -->
          <div id="notes-grid" class="min-h-[400px]">
             <NotesList
              :notes="notes"
              :loading="loading"
              @delete-note="handleDeleteNote"
              @edit-note="handleEditNote"
            />
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="flex items-center justify-center gap-4 pt-4">
            <button
              @click="changePage(page - 1)"
              :disabled="page === 1"
              class="group flex items-center gap-2 px-4 py-2 rounded-full bg-white border border-slate-200 text-sm font-medium text-slate-600 hover:border-primary-200 hover:text-primary-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              <svg class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              Prev
            </button>
            
            <div class="flex items-center gap-1">
              <span class="text-sm font-medium text-slate-400">Page</span>
              <span class="text-sm font-bold text-slate-900">{{ page }}</span>
              <span class="text-sm font-medium text-slate-400">of {{ totalPages }}</span>
            </div>

            <button
              @click="changePage(page + 1)"
              :disabled="page === totalPages"
              class="group flex items-center gap-2 px-4 py-2 rounded-full bg-white border border-slate-200 text-sm font-medium text-slate-600 hover:border-primary-200 hover:text-primary-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              Next
              <svg class="w-4 h-4 group-hover:translate-x-0.5 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>

        </main>
      </div>
    </div>
  </div>
</template>
