<script setup>
const props = defineProps({
  note: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['delete', 'edit'])

function onDelete() {
  emit('delete', props.note.id)
}

function formatDate(dateString) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}
</script>

<template>
  <li
    class="group relative flex flex-col justify-between h-full bg-white rounded-2xl border border-slate-200/60 p-5 shadow-sm hover:shadow-md hover:border-primary-200 transition-all duration-300 hover:-translate-y-1"
  >
    <div class="space-y-3">
      <div class="flex items-start justify-between gap-2">
        <h3 class="font-display font-bold text-lg text-slate-900 leading-tight line-clamp-2 group-hover:text-primary-600 transition-colors">
          {{ note.title }}
        </h3>
      </div>
      
      <p
        v-if="note.content"
        class="text-sm text-slate-500 leading-relaxed line-clamp-3"
      >
        {{ note.content }}
      </p>
      <p
        v-else
        class="text-sm italic text-slate-400"
      >
        No content provided.
      </p>
    </div>

    <div class="mt-5 pt-4 border-t border-slate-100 flex items-center justify-between">
      <span class="text-xs font-medium text-slate-400">
        {{ formatDate(note.created_at) }}
      </span>

      <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <button
          @click="$emit('edit', note)"
          class="p-2 rounded-lg text-slate-400 hover:text-primary-600 hover:bg-primary-50 transition-colors"
          title="Edit"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button
          @click="onDelete"
          class="p-2 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
          title="Delete"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
  </li>
</template>
