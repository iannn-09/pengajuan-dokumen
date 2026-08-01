<template>
  <div class="tiptap-editor-wrapper" v-if="editor">
    <!-- Rich Text Toolbar -->
    <div class="editor-toolbar">
      <button 
        type="button"
        class="toolbar-btn" 
        :class="{ active: editor.isActive('bold') }" 
        @click="editor.chain().focus().toggleBold().run()"
        title="Tebal (Bold)"
      >
        <strong>B</strong>
      </button>

      <button 
        type="button"
        class="toolbar-btn" 
        :class="{ active: editor.isActive('italic') }" 
        @click="editor.chain().focus().toggleItalic().run()"
        title="Miring (Italic)"
      >
        <em>I</em>
      </button>

      <button 
        type="button"
        class="toolbar-btn" 
        :class="{ active: editor.isActive('strike') }" 
        @click="editor.chain().focus().toggleStrike().run()"
        title="Coret (Strikethrough)"
      >
        <s>S</s>
      </button>

      <div class="toolbar-divider"></div>

      <button 
        type="button"
        class="toolbar-btn" 
        :class="{ active: editor.isActive('bulletList') }" 
        @click="editor.chain().focus().toggleBulletList().run()"
        title="Daftar Simbol (Bullet List)"
      >
        • List
      </button>

      <button 
        type="button"
        class="toolbar-btn" 
        :class="{ active: editor.isActive('orderedList') }" 
        @click="editor.chain().focus().toggleOrderedList().run()"
        title="Daftar Angka (Numbered List)"
      >
        1. List
      </button>

      <div class="toolbar-divider"></div>

      <button 
        type="button"
        class="toolbar-btn" 
        :class="{ active: editor.isActive('heading', { level: 3 }) }" 
        @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
        title="Judul Sub-Bab (H3)"
      >
        H3
      </button>

      <div class="toolbar-divider"></div>

      <button 
        type="button"
        class="toolbar-btn" 
        @click="editor.chain().focus().unsetAllMarks().clearNodes().run()"
        title="Hapus Format"
      >
        🧹
      </button>

      <button 
        type="button"
        class="toolbar-btn" 
        @click="editor.chain().focus().undo().run()"
        :disabled="!editor.can().undo()"
        title="Urungkan (Undo)"
      >
        ↺
      </button>

      <button 
        type="button"
        class="toolbar-btn" 
        @click="editor.chain().focus().redo().run()"
        :disabled="!editor.can().redo()"
        title="Ulangi (Redo)"
      >
        ↻
      </button>
    </div>

    <!-- Editor Content Body -->
    <EditorContent :editor="editor" class="editor-body" />
  </div>
</template>

<script setup>
import { watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Tuliskan isian teks di sini...'
  }
})

const emit = defineEmits(['update:modelValue'])

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit
  ],
  onUpdate: () => {
    emit('update:modelValue', editor.value.getHTML())
  }
})

watch(() => props.modelValue, (newValue) => {
  const isSame = editor.value.getHTML() === newValue
  if (isSame) return
  editor.value.commands.setContent(newValue, false)
})

onBeforeUnmount(() => {
  if (editor.value) {
    editor.value.destroy()
  }
})
</script>

<style scoped>
.tiptap-editor-wrapper {
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.5rem 0.75rem;
  background: rgba(15, 23, 42, 0.6);
  border-bottom: 1px solid var(--border-color);
  flex-wrap: wrap;
}

.toolbar-btn {
  padding: 0.3rem 0.6rem;
  font-size: 0.82rem;
  background: transparent;
  border: 1px solid transparent;
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s ease;
  min-width: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.toolbar-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-main);
}

.toolbar-btn.active {
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8;
  border-color: rgba(99, 102, 241, 0.4);
}

.toolbar-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.toolbar-divider {
  width: 1px;
  height: 16px;
  background: var(--border-color);
  margin: 0 0.2rem;
}

.editor-body {
  padding: 0.85rem 1rem;
  min-height: 140px;
  max-height: 350px;
  overflow-y: auto;
  color: var(--text-main);
  font-size: 0.9rem;
  line-height: 1.6;
}

:deep(.ProseMirror) {
  outline: none;
  min-height: 120px;
}

:deep(.ProseMirror p) {
  margin-bottom: 0.5rem;
}

:deep(.ProseMirror ul), :deep(.ProseMirror ol) {
  padding-left: 1.5rem;
  margin-bottom: 0.5rem;
}

:deep(.ProseMirror h3) {
  font-size: 1.05rem;
  font-weight: 700;
  margin-top: 0.5rem;
  margin-bottom: 0.35rem;
  color: var(--text-main);
}
</style>
