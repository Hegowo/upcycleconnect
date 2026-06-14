<template>
  <div class="tiptap-editor-wrap border border-[#e5e7eb] rounded-xl overflow-hidden focus-within:ring-2 focus-within:ring-[#006d35]/30 focus-within:border-[#006d35]/50 transition">

    <div class="tiptap-toolbar flex flex-wrap items-center gap-0.5 p-2 bg-[#f8fafc] border-b border-[#e5e7eb]">

      <select @change="setHeading($event.target.value)" class="h-7 px-2 text-xs border border-[#e5e7eb] rounded-lg bg-white text-[#374151] focus:outline-none cursor-pointer mr-1">
        <option value="p">Paragraphe</option>
        <option value="1">Titre 1</option>
        <option value="2">Titre 2</option>
        <option value="3">Titre 3</option>
      </select>

      <div class="h-5 w-px bg-[#e5e7eb] mx-1" />

      <ToolBtn @click="editor.chain().focus().toggleBold().run()" :active="editor?.isActive('bold')" title="Gras">
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M5 4a1 1 0 0 1 1-1h4.5a4 4 0 0 1 2.828 6.828A4 4 0 0 1 10.5 17H6a1 1 0 0 1-1-1V4zm2 1v4h3.5a2 2 0 0 0 0-4H7zm0 6v4h3.5a2 2 0 0 0 0-4H7z"/></svg>
      </ToolBtn>
      <ToolBtn @click="editor.chain().focus().toggleItalic().run()" :active="editor?.isActive('italic')" title="Italique">
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M8 4a1 1 0 0 1 1-1h4a1 1 0 1 1 0 2h-1.28l-2.44 10H11a1 1 0 1 1 0 2H7a1 1 0 1 1 0-2h1.28L10.72 5H9a1 1 0 0 1-1-1z"/></svg>
      </ToolBtn>
      <ToolBtn @click="editor.chain().focus().toggleUnderline().run()" :active="editor?.isActive('underline')" title="Souligné">
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M5 15a1 1 0 0 0 0 2h10a1 1 0 0 0 0-2H5zM6 4a1 1 0 0 1 1-1h6a1 1 0 1 1 0 2V9a3 3 0 0 1-6 0V5a1 1 0 0 1-1-1z"/></svg>
      </ToolBtn>
      <ToolBtn @click="editor.chain().focus().toggleStrike().run()" :active="editor?.isActive('strike')" title="Barré">
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M3 10a1 1 0 0 1 1-1h12a1 1 0 1 1 0 2H4a1 1 0 0 1-1-1zM7 4a1 1 0 0 1 1-1h4a1 1 0 1 1 0 2H8a1 1 0 0 1-1-1zM7.293 15.293a1 1 0 0 1 1.414 0l1.293 1.293 1.293-1.293a1 1 0 0 1 1.414 1.414l-2 2a1 1 0 0 1-1.414 0l-2-2a1 1 0 0 1 0-1.414z"/></svg>
      </ToolBtn>

      <div class="h-5 w-px bg-[#e5e7eb] mx-1" />

      <ToolBtn @click="editor.chain().focus().toggleBulletList().run()" :active="editor?.isActive('bulletList')" title="Liste à puces">
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M4 5a1 1 0 1 1-2 0 1 1 0 0 1 2 0zM6 5a1 1 0 0 1 1-1h9a1 1 0 1 1 0 2H7a1 1 0 0 1-1-1zM4 10a1 1 0 1 1-2 0 1 1 0 0 1 2 0zM6 10a1 1 0 0 1 1-1h9a1 1 0 1 1 0 2H7a1 1 0 0 1-1-1zM4 15a1 1 0 1 1-2 0 1 1 0 0 1 2 0zM6 15a1 1 0 0 1 1-1h9a1 1 0 1 1 0 2H7a1 1 0 0 1-1-1z"/></svg>
      </ToolBtn>
      <ToolBtn @click="editor.chain().focus().toggleOrderedList().run()" :active="editor?.isActive('orderedList')" title="Liste numérotée">
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M3 4a1 1 0 0 1 2 0v3H4a1 1 0 0 1 0-2H3V4zM3 9a1 1 0 0 0 0 2h.5v1H3a1 1 0 0 0 0 2h1a1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1H3zM6 5a1 1 0 0 1 1-1h9a1 1 0 1 1 0 2H7a1 1 0 0 1-1-1zM6 10a1 1 0 0 1 1-1h9a1 1 0 1 1 0 2H7a1 1 0 0 1-1-1zM6 15a1 1 0 0 1 1-1h9a1 1 0 1 1 0 2H7a1 1 0 0 1-1-1z"/></svg>
      </ToolBtn>
      <ToolBtn @click="editor.chain().focus().toggleBlockquote().run()" :active="editor?.isActive('blockquote')" title="Citation">
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M6.5 3a3.5 3.5 0 0 1 3.465 3.956A3.502 3.502 0 0 1 8 10.5a3.5 3.5 0 0 1-3-1.708V10a6 6 0 0 0 6 6 1 1 0 1 1 0 2 8 8 0 0 1-8-8V6.5A3.5 3.5 0 0 1 6.5 3zm8 0a3.5 3.5 0 0 1 3.465 3.956A3.502 3.502 0 0 1 16 10.5a3.5 3.5 0 0 1-3-1.708V10a6 6 0 0 0 6 6 1 1 0 1 1 0 2 8 8 0 0 1-8-8V6.5A3.5 3.5 0 0 1 14.5 3z"/></svg>
      </ToolBtn>

      <div class="h-5 w-px bg-[#e5e7eb] mx-1" />

      <div class="relative flex items-center" title="Couleur du texte">
        <input type="color" :value="currentColor" @input="e => editor.chain().focus().setColor(e.target.value).run()"
          class="w-7 h-7 rounded-lg cursor-pointer border border-[#e5e7eb] p-0.5 bg-white" />
      </div>

      <div class="relative flex items-center" title="Surlignage">
        <input type="color" value="#fef08a" @input="e => editor.chain().focus().toggleHighlight({ color: e.target.value }).run()"
          class="w-7 h-7 rounded-lg cursor-pointer border border-[#e5e7eb] p-0.5 bg-[#fef08a]" />
      </div>

      <div class="h-5 w-px bg-[#e5e7eb] mx-1" />

      <ToolBtn @click="addLink" :active="editor?.isActive('link')" title="Lien">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M13.828 10.172a4 4 0 0 0-5.656 0l-4 4a4 4 0 1 0 5.656 5.656l1.102-1.101m-.758-4.899a4 4 0 0 0 5.656 0l4-4a4 4 0 0 0-5.656-5.656l-1.1 1.1"/></svg>
      </ToolBtn>

      <ToolBtn @click="triggerImageInput" title="Image">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909M3 21h18a.75.75 0 0 0 .75-.75V5.25A.75.75 0 0 0 21 3H3a.75.75 0 0 0-.75.75v15c0 .414.336.75.75.75ZM9.75 9.75a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z"/></svg>
      </ToolBtn>
      <input ref="imgInput" type="file" accept="image/*" class="hidden" @change="uploadImage" />

      <div class="h-5 w-px bg-[#e5e7eb] mx-1" />

      <ToolBtn @click="editor.chain().focus().undo().run()" :disabled="!editor?.can().undo()" title="Annuler">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9 15 3 9m0 0 6-6M3 9h12a6 6 0 0 1 0 12h-3"/></svg>
      </ToolBtn>
      <ToolBtn @click="editor.chain().focus().redo().run()" :disabled="!editor?.can().redo()" title="Rétablir">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="m15 15 6-6m0 0-6-6m6 6H9a6 6 0 0 0 0 12h3"/></svg>
      </ToolBtn>

      <span v-if="uploading" class="ml-2 text-xs text-[#94a3b8] flex items-center gap-1">
        <div class="w-3 h-3 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin" />
        Upload...
      </span>
    </div>

    <EditorContent :editor="editor" class="tiptap-content" />
  </div>
</template>

<script setup>
import { ref, computed, watch, onBeforeUnmount, defineComponent, h } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import { Color } from '@tiptap/extension-color'
import { TextStyle } from '@tiptap/extension-text-style'
import Highlight from '@tiptap/extension-highlight'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { useUserAuthStore } from '@/stores/userAuth'

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: 'Écrivez votre message...' },
  minHeight: { type: String, default: '180px' },
})
const emit = defineEmits(['update:modelValue'])

const userAuth = useUserAuthStore()
const imgInput = ref(null)
const uploading = ref(false)

const ToolBtn = defineComponent({
  props: { active: Boolean, disabled: Boolean },
  setup(props, { slots, attrs }) {
    return () => h('button', {
      type: 'button',
      disabled: props.disabled,
      ...attrs,
      class: [
        'w-7 h-7 flex items-center justify-center rounded-lg text-sm transition',
        props.active ? 'bg-[#006d35] text-white' : 'text-[#374151] hover:bg-[#e5e7eb]',
        props.disabled ? 'opacity-30 cursor-not-allowed' : 'cursor-pointer',
      ].join(' '),
    }, slots.default?.())
  },
})

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    Underline,
    TextStyle,
    Color,
    Highlight.configure({ multicolor: true }),
    Image.configure({ allowBase64: true }),
    Link.configure({ openOnClick: false }),
    Placeholder.configure({ placeholder: props.placeholder }),
  ],
  editorProps: {
    attributes: {
      class: 'tiptap-inner focus:outline-none',
      style: `min-height: ${props.minHeight}`,
    },
  },
  onUpdate({ editor }) {
    emit('update:modelValue', editor.getHTML())
  },
})

const currentColor = computed(() => editor.value?.getAttributes('textStyle').color || '#000000')

watch(() => props.modelValue, (val) => {
  if (editor.value && editor.value.getHTML() !== val) {
    editor.value.commands.setContent(val, false)
  }
})

function setHeading(val) {
  if (val === 'p') {
    editor.value?.chain().focus().setParagraph().run()
  } else {
    editor.value?.chain().focus().setHeading({ level: parseInt(val) }).run()
  }
}

function addLink() {
  const url = window.prompt('URL du lien :')
  if (url) {
    editor.value?.chain().focus().setLink({ href: url }).run()
  } else if (url === '') {
    editor.value?.chain().focus().unsetLink().run()
  }
}

function triggerImageInput() {
  imgInput.value?.click()
}

async function uploadImage(e) {
  const file = e.target.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    const form = new FormData()
    form.append('file', file)
    const res = await fetch('/api/v1/forum/media', {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
      body: form,
    })
    if (res.ok) {
      const data = await res.json()
      editor.value?.chain().focus().setImage({ src: data.url }).run()
    }
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

onBeforeUnmount(() => editor.value?.destroy())
</script>

<style>
.tiptap-content .tiptap-inner {
  padding: 12px 16px;
  line-height: 1.7;
  color: #1e293b;
}

.tiptap-content .tiptap-inner p { margin-bottom: 0.75em; }
.tiptap-content .tiptap-inner p:last-child { margin-bottom: 0; }

.tiptap-content .tiptap-inner h1 { font-size: 1.5rem; font-weight: 800; margin-bottom: 0.5em; color: #001d32; }
.tiptap-content .tiptap-inner h2 { font-size: 1.25rem; font-weight: 700; margin-bottom: 0.5em; color: #001d32; }
.tiptap-content .tiptap-inner h3 { font-size: 1.1rem; font-weight: 700; margin-bottom: 0.5em; color: #001d32; }

.tiptap-content .tiptap-inner ul { list-style: disc; padding-left: 1.5em; margin-bottom: 0.75em; }
.tiptap-content .tiptap-inner ol { list-style: decimal; padding-left: 1.5em; margin-bottom: 0.75em; }
.tiptap-content .tiptap-inner li { margin-bottom: 0.25em; }

.tiptap-content .tiptap-inner blockquote {
  border-left: 4px solid #006d35;
  padding-left: 1em;
  color: #64748b;
  font-style: italic;
  margin: 0.75em 0;
}

.tiptap-content .tiptap-inner code {
  background: #f1f5f9;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 0.875em;
  font-family: monospace;
}

.tiptap-content .tiptap-inner pre {
  background: #1e293b;
  color: #e2e8f0;
  border-radius: 8px;
  padding: 12px 16px;
  overflow-x: auto;
  margin-bottom: 0.75em;
}

.tiptap-content .tiptap-inner pre code {
  background: none;
  padding: 0;
  color: inherit;
}

.tiptap-content .tiptap-inner img {
  max-width: 100%;
  border-radius: 8px;
  margin: 0.5em 0;
}

.tiptap-content .tiptap-inner a {
  color: #006d35;
  text-decoration: underline;
}

.tiptap-content .tiptap-inner hr {
  border: none;
  border-top: 2px solid #e5e7eb;
  margin: 1em 0;
}

.tiptap-content .tiptap-inner mark {
  border-radius: 3px;
  padding: 0 2px;
}

.tiptap-content .tiptap-inner p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: #94a3b8;
  pointer-events: none;
  height: 0;
}
</style>
