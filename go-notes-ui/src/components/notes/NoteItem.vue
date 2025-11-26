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
    class="group relative flex flex-col justify-between gap-4 rounded-2xl bg-white p-5 shadow-sm border border-slate-100 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-slate-200/50"
  >
    <!-- Top Accent (Random color simulation based on ID parity for now, or just a gradient) -->
    <div 
      class="absolute top-0 left-0 right-0 h-1.5 rounded-t-2xl bg-gradient-to-r from-primary-400 to-indigo-400 opacity-0 group-hover:opacity-100 transition-opacity"
    ></div>

    <div class="space-y-2">
      <div class="flex items-start justify-between gap-2">
        <h3 class="font-bold text-lg text-slate-800 leading-tight group-hover:text-primary-600 transition-colors">
          {{ note.title }}
        </h3>
      </div>
      
      <p
        v-if="note.content"
        class="text-sm leading-relaxed text-slate-600 line-clamp-4"
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

    <div class="flex items-center justify-between pt-2 border-t border-slate-50 mt-2">
      <span class="text-[10px] font-medium text-slate-400 uppercase tracking-wider">
        {{ formatDate(note.created_at) }}
      </span>

      <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all transform translate-y-2 group-hover:translate-y-0">
        <button
          class="p-2 rounded-lg text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 transition-colors"
          @click="$emit('edit', note)"
          title="Edit"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button
          class="p-2 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
          @click="onDelete"
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
