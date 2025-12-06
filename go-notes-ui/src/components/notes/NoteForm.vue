<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  noteToEdit: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['create-note', 'update-note', 'cancel-edit'])

const title = ref('')
const content = ref('')

watch(
  () => props.noteToEdit,
  (newVal) => {
    if (newVal) {
      title.value = newVal.title
      content.value = newVal.content
    } else {
      title.value = ''
      content.value = ''
    }
  },
  { immediate: true }
)

function onSubmit() {
  if (props.noteToEdit) {
    emit('update-note', {
      id: props.noteToEdit.id,
      title: title.value,
      content: content.value,
    })
  } else {
    emit('create-note', {
      title: title.value,
      content: content.value,
    })
    title.value = ''
    content.value = ''
  }
}

function onCancel() {
  emit('cancel-edit')
}
</script>

<template>
  <div class="group relative">
    <!-- Glow effect behind -->
    <div class="absolute -inset-1 bg-gradient-to-r from-primary-200 to-indigo-200 rounded-2xl blur opacity-25 group-hover:opacity-50 transition duration-500"></div>
    
    <section class="relative bg-white dark:bg-slate-800 rounded-2xl shadow-xl shadow-slate-200/50 dark:shadow-none border border-slate-100 dark:border-slate-700 overflow-hidden transition-colors duration-300">
      <div class="p-1">
        <form @submit.prevent="onSubmit">
          <!-- Title Input -->
          <div class="relative">
            <input
              v-model="title"
              type="text"
              placeholder="What's on your mind?"
              class="w-full px-5 py-4 text-lg font-semibold text-slate-800 dark:text-slate-100 placeholder:text-slate-400 dark:placeholder:text-slate-500 bg-transparent border-none outline-none focus:ring-0 transition-colors"
            />
          </div>

          <!-- Content Textarea -->
          <div class="relative px-5 pb-2">
            <textarea
              v-model="content"
              rows="2"
              placeholder="Add some details..."
              class="w-full text-slate-600 dark:text-slate-300 placeholder:text-slate-400 dark:placeholder:text-slate-500 bg-transparent border-none outline-none focus:ring-0 resize-none text-sm leading-relaxed transition-colors"
            ></textarea>
          </div>

          <!-- Footer / Actions -->
          <div class="flex items-center justify-between px-4 py-3 bg-slate-50/50 dark:bg-slate-900/50 border-t border-slate-100 dark:border-slate-700 transition-colors">
            <div class="flex items-center gap-2">
              <span v-if="noteToEdit" class="text-xs font-medium text-primary-600 dark:text-primary-400 bg-primary-50 dark:bg-primary-900/30 px-2 py-1 rounded-md">
                Editing #{{ noteToEdit.id }}
              </span>
            </div>

            <div class="flex items-center gap-3">
              <button
                v-if="noteToEdit"
                type="button"
                @click="onCancel"
                class="text-sm font-medium text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
              >
                Cancel
              </button>
              
              <button
                type="submit"
                :disabled="!title.trim()"
                class="inline-flex items-center gap-2 px-5 py-2 bg-slate-900 dark:bg-white hover:bg-slate-800 dark:hover:bg-slate-100 text-white dark:text-slate-900 text-sm font-medium rounded-xl transition-all shadow-lg shadow-slate-900/20 dark:shadow-none hover:shadow-slate-900/30 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span>{{ noteToEdit ? 'Update Note' : 'Create Note' }}</span>
                <svg v-if="!noteToEdit" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </button>
            </div>
          </div>
        </form>
      </div>
    </section>
  </div>
</template>
