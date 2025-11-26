<script setup>
import { ref, onMounted } from 'vue'
import { getNotes, createNote, deleteNote, updateNote } from '../../api/notes'
import NoteForm from './NoteForm.vue'
import NotesList from './NotesList.vue'

const notes = ref([])
const loading = ref(false)
const error = ref('')
const editingNote = ref(null)

async function loadNotes() {
  loading.value = true
  error.value = ''
  try {
    notes.value = await getNotes()
  } catch (err) {
    error.value = err.message || 'Terjadi kesalahan'
  } finally {
    loading.value = false
  }
}

async function handleCreateNote({ title, content }) {
  if (!title.trim()) {
    error.value = 'Title wajib diisi'
    return
  }

  try {
    error.value = ''
    const created = await createNote({ title, content })
    notes.value.push(created)
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
    notes.value = notes.value.filter((n) => n.id !== id)
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
    const updated = await updateNote(id, { title, content })
    const index = notes.value.findIndex((n) => n.id === id)
    if (index !== -1) {
      notes.value[index] = updated
    }
    editingNote.value = null
  } catch (err) {
    error.value = err.message || 'Gagal mengupdate note'
  }
}

onMounted(loadNotes)
</script>

<template>
  <div class="min-h-screen bg-slate-50 font-sans selection:bg-primary-100 selection:text-primary-900">
    <!-- Top Decoration -->
    <div class="fixed inset-0 pointer-events-none overflow-hidden">
      <div class="absolute -top-[20%] -left-[10%] w-[50%] h-[50%] rounded-full bg-primary-200/20 blur-3xl"></div>
      <div class="absolute top-[10%] -right-[10%] w-[40%] h-[40%] rounded-full bg-indigo-200/20 blur-3xl"></div>
    </div>

    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Header -->
      <header class="flex flex-col items-center text-center mb-12 space-y-4">
        <div class="inline-flex items-center justify-center p-2 bg-white rounded-2xl shadow-sm border border-slate-100 mb-2">
          <span class="text-2xl">⚡</span>
        </div>
        <h1 class="text-4xl sm:text-5xl font-bold tracking-tight text-slate-900">
          Super <span class="text-primary-600">Notes</span>
        </h1>
        <p class="max-w-lg text-lg text-slate-600">
          Capture your ideas instantly. Simple, fast, and beautiful.
        </p>
        
        <!-- Status Indicator -->
        <div class="flex items-center gap-2 mt-2">
          <span class="relative flex h-2.5 w-2.5">
            <span v-if="loading" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-400 opacity-75"></span>
            <span :class="loading ? 'bg-amber-500' : 'bg-emerald-500'" class="relative inline-flex rounded-full h-2.5 w-2.5"></span>
          </span>
          <span class="text-xs font-medium text-slate-500 uppercase tracking-wider">
            {{ loading ? 'Syncing...' : 'All Systems Go' }}
          </span>
        </div>
      </header>

      <!-- Error Toast -->
      <transition
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="transform -translate-y-2 opacity-0"
        enter-to-class="transform translate-y-0 opacity-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="transform translate-y-0 opacity-100"
        leave-to-class="transform -translate-y-2 opacity-0"
      >
        <div v-if="error" class="max-w-md mx-auto mb-8 p-4 bg-rose-50 border border-rose-100 rounded-xl flex items-center gap-3 text-rose-700 shadow-sm">
          <svg class="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-sm font-medium">{{ error }}</span>
        </div>
      </transition>

      <!-- Main Content -->
      <main class="max-w-3xl mx-auto space-y-12">
        <!-- Input Area -->
        <section class="relative z-10">
          <NoteForm
            :note-to-edit="editingNote"
            @create-note="handleCreateNote"
            @update-note="handleUpdateNote"
            @cancel-edit="handleCancelEdit"
          />
        </section>

        <!-- Notes Grid -->
        <section>
          <div class="flex items-center justify-between mb-6 px-2">
            <h2 class="text-xl font-bold text-slate-800 flex items-center gap-2">
              Your Notes
              <span class="px-2.5 py-0.5 rounded-full bg-slate-100 text-slate-600 text-xs font-bold">
                {{ notes.length }}
              </span>
            </h2>
          </div>

          <NotesList
            :notes="notes"
            :loading="loading"
            @delete-note="handleDeleteNote"
            @edit-note="handleEditNote"
          />
        </section>
      </main>
    </div>
  </div>
</template>
