<script setup>
import { ref, onMounted } from 'vue'
import { getNotes, createNote, deleteNote } from '../../api/notes'
import NoteForm from './NoteForm.vue'
import NotesList from './NotesList.vue'

const notes = ref([])
const loading = ref(false)
const error = ref('')

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
  } catch (err) {
    error.value = err.message || 'Gagal menghapus note'
  }
}

onMounted(loadNotes)
</script>

<template>
  <div
    class="min-h-screen bg-gradient-to-br from-slate-50 via-slate-100 to-slate-200 px-4 py-8 flex justify-center"
  >
    <div class="w-full max-w-5xl">
      <!-- Header -->
      <header class="flex items-start justify-between gap-4">
        <div>
          <h1 class="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-900">
            Notes API Demo
          </h1>
          <p class="mt-1 text-sm text-slate-600">
            Kelola ide dan tugas Anda dengan cepat.
          </p>
          <p class="text-xs text-slate-400">
            Go REST API + Vue 3 + Tailwind
          </p>
        </div>

        <span
          v-if="loading"
          class="mt-1 inline-flex items-center gap-1 rounded-full bg-amber-50 px-3 py-1 text-[11px] font-medium text-amber-700 ring-1 ring-amber-200"
        >
          <span class="h-1.5 w-1.5 rounded-full bg-amber-500 animate-pulse" />
          Loading…
        </span>
        <span
          v-else
          class="mt-1 inline-flex items-center gap-1 rounded-full bg-emerald-50 px-3 py-1 text-[11px] font-medium text-emerald-700 ring-1 ring-emerald-200"
        >
          <span class="h-1.5 w-1.5 rounded-full bg-emerald-500" />
          Ready
        </span>
      </header>

      <!-- Error -->
      <section
        v-if="error"
        class="mt-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-2.5 text-xs sm:text-sm text-rose-800"
      >
        {{ error }}
      </section>

      <!-- Main content: 2 kolom -->
      <div
        class="mt-6 grid gap-6 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,1.25fr)] items-start"
      >
        <!-- Form kiri -->
        <NoteForm @create-note="handleCreateNote" />

        <!-- List kanan -->
        <section
          class="rounded-2xl border border-slate-200 bg-white/90 shadow-lg shadow-slate-200/70 px-4 py-4 sm:px-5 sm:py-5"
        >
          <div class="mb-3 flex items-center justify-between">
            <h2 class="text-sm font-semibold text-slate-900">
              Daftar Notes
            </h2>
            <span class="text-[11px] text-slate-500">
              {{ notes.length }} item
            </span>
          </div>

          <NotesList
            :notes="notes"
            :loading="loading"
            @delete-note="handleDeleteNote"
          />
        </section>
      </div>
    </div>
  </div>
</template>
