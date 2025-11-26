<script setup>
const props = defineProps({
  note: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['delete'])

function onDelete() {
  emit('delete', props.note.id)
}
</script>

<template>
  <li
    class="group flex items-start justify-between gap-3 rounded-xl border border-slate-200 bg-slate-50 px-3.5 py-3 text-sm text-slate-900 shadow-sm shadow-slate-200/80 transition hover:-translate-y-0.5 hover:border-indigo-300 hover:bg-white hover:shadow-md hover:shadow-indigo-100"
  >
    <div class="space-y-1">
      <h3 class="font-semibold text-slate-900">
        {{ note.title }}
      </h3>
      <p
        v-if="note.content"
        class="text-xs sm:text-[13px] leading-relaxed text-slate-600"
      >
        {{ note.content }}
      </p>
      <p
        v-else
        class="text-xs italic text-slate-400"
      >
        Tidak ada konten.
      </p>
    </div>

    <button
      class="mt-0.5 inline-flex h-7 w-7 items-center justify-center rounded-full border border-rose-200 bg-rose-50 text-[11px] font-semibold text-rose-600 opacity-0 shadow-sm shadow-rose-100 transition group-hover:opacity-100 hover:bg-rose-500 hover:border-rose-500 hover:text-white active:scale-95"
      @click="onDelete"
      title="Hapus"
    >
      ✕
    </button>
  </li>
</template>
