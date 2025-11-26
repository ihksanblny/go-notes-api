<script setup>
import NoteItem from './NoteItem.vue'

const props = defineProps({
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
    <p
      v-if="!loading && notes.length === 0"
      class="mt-1.5 text-xs sm:text-sm text-slate-400"
    >
      Belum ada note. Tambah satu di form di sebelah kiri.
    </p>

    <ul
      v-else
      class="mt-1.5 space-y-2.5"
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
