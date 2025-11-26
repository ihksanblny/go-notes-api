<script setup>
import NoteItem from './NoteItem.vue'

defineProps({
  notes: {
    type: Array,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['delete-note', 'edit-note'])

function handleDelete(id) {
  emit('delete-note', id)
}
</script>

<template>
  <div>
    <div
      v-if="!loading && notes.length === 0"
      class="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed border-slate-200 rounded-2xl bg-slate-50/50"
    >
      <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mb-4 text-3xl">
        📝
      </div>
      <h3 class="text-slate-900 font-medium mb-1">No notes yet</h3>
      <p class="text-slate-500 text-sm max-w-xs">
        Your ideas are safe here. Start by creating your first note above.
      </p>
    </div>

    <ul
      v-else
      class="grid grid-cols-1 sm:grid-cols-2 gap-4"
    >
      <NoteItem
        v-for="note in notes"
        :key="note.id"
        :note="note"
        @delete="handleDelete"
        @edit="$emit('edit-note', $event)"
      />
    </ul>
  </div>
</template>
