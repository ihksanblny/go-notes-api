<script setup>
import { ref } from 'vue'

const emit = defineEmits(['create-note'])

const title = ref('')
const content = ref('')

function onSubmit() {
  emit('create-note', {
    title: title.value,
    content: content.value,
  })

  title.value = ''
  content.value = ''
}
</script>

<template>
  <section
    class="rounded-2xl border border-slate-200 bg-white/95 shadow-lg shadow-slate-200/70 px-4 py-4 sm:px-5 sm:py-5"
  >
    <h2 class="mb-3 text-sm font-semibold text-slate-900">
      Buat Note
    </h2>

    <form @submit.prevent="onSubmit" class="space-y-3">
      <div class="space-y-1.5">
        <label class="block text-[11px] font-semibold tracking-wide text-slate-500 uppercase">
          Title
        </label>
        <input
          v-model="title"
          type="text"
          placeholder="Misal: Clean Code, Keamanan…"
          class="block w-full rounded-xl border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-900 placeholder:text-slate-400 outline-none transition focus:border-indigo-400 focus:bg-white focus:ring-2 focus:ring-indigo-100"
        />
      </div>

      <div class="space-y-1.5">
        <label class="block text-[11px] font-semibold tracking-wide text-slate-500 uppercase">
          Content
        </label>
        <textarea
          v-model="content"
          rows="3"
          placeholder="Deskripsi singkat (opsional)…"
          class="block w-full resize-none rounded-xl border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-900 placeholder:text-slate-400 outline-none transition focus:border-indigo-400 focus:bg-white focus:ring-2 focus:ring-indigo-100"
        ></textarea>
      </div>

      <div class="flex justify-between items-center pt-1">
        <span class="text-[11px] text-slate-400">
          Tip: pakai kata kerja aktif untuk judul.
        </span>
        <button
          type="submit"
          class="inline-flex items-center gap-1 rounded-full bg-gradient-to-r from-indigo-500 to-sky-500 px-4 py-1.5 text-xs font-semibold text-white shadow-md shadow-indigo-300/70 transition hover:shadow-lg hover:shadow-indigo-300/90 hover:brightness-110 active:scale-[0.97]"
        >
          <span class="text-base leading-none">+</span>
          Tambah Note
        </button>
      </div>
    </form>
  </section>
</template>
