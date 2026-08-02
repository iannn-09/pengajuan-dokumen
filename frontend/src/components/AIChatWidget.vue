<template>
  <div class="ai-widget-wrapper">
    <!-- Floating Action Button (FAB) -->
    <button 
      class="ai-fab-btn" 
      @click="isOpen = !isOpen" 
      :title="isOpen ? 'Tutup Asisten AI' : 'Buka Asisten AI Kelayakan'"
    >
      <span class="ai-fab-icon">🤖</span>
      <span class="ai-fab-text">Asisten AI</span>
      <span class="pulse-dot"></span>
    </button>

    <!-- Floating Chat Window Dialog -->
    <div v-if="isOpen" class="ai-chat-window glass-card">
      <div class="chat-header">
        <div class="header-left">
          <div class="ai-avatar">🤖</div>
          <div>
            <h4 class="chat-title">Asisten AI Kelayakan</h4>
            <span class="chat-status"><span class="green-dot"></span> Online & Siap Membantu</span>
          </div>
        </div>
        <button class="close-chat-btn" @click="isOpen = false">✕</button>
      </div>

      <div class="chat-body" ref="chatBodyRef">
        <!-- Welcoming Message & Quick Prompt Chips -->
        <div class="chat-msg msg-ai">
          <div class="msg-avatar">🤖</div>
          <div class="msg-bubble">
            <p>Halo <strong>{{ auth.userName }}</strong>! 👋 Saya **Asisten AI Kelayakan**.</p>
            <p class="mt-1">Ada yang bisa saya bantu terkait pengajuan dokumen, aturan berkas, atau status permohonan hari ini?</p>
          </div>
        </div>

        <!-- Dynamic Chat Messages History -->
        <div 
          v-for="(msg, index) in messages" 
          :key="index" 
          class="chat-msg" 
          :class="msg.role === 'user' ? 'msg-user' : 'msg-ai'"
        >
          <div v-if="msg.role === 'ai'" class="msg-avatar">🤖</div>
          <div class="msg-bubble" v-html="formatMessage(msg.content)"></div>
        </div>

        <!-- Typing Indicator -->
        <div v-if="loading" class="chat-msg msg-ai">
          <div class="msg-avatar">🤖</div>
          <div class="msg-bubble typing-bubble">
            <span class="dot-typing"></span>
            <span class="dot-typing"></span>
            <span class="dot-typing"></span>
          </div>
        </div>
      </div>

      <!-- Quick Suggestion Chips Bar -->
      <div class="chips-scroll">
        <button 
          v-for="chip in quickChips" 
          :key="chip" 
          class="chip-btn" 
          @click="sendQuickPrompt(chip)"
          :disabled="loading"
        >
          {{ chip }}
        </button>
      </div>

      <!-- Footer Message Input -->
      <div class="chat-footer">
        <input 
          v-model="inputMessage" 
          type="text" 
          class="chat-input" 
          placeholder="Ketik pertanyaan Anda di sini..." 
          @keyup.enter="sendMessage"
          :disabled="loading"
        />
        <button class="chat-send-btn" @click="sendMessage" :disabled="loading || !inputMessage.trim()">
          🚀
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'

const auth = useAuthStore()
const isOpen = ref(false)
const loading = ref(false)
const inputMessage = ref('')
const chatBodyRef = ref(null)

const messages = ref([])

const quickChips = [
  '📌 Syarat Berkas Kelayakan',
  '📊 Penjelasan Status Permohonan',
  '📝 Cara Perbaiki Revisi',
  '⏱️ Berapa Lama Proses Verifikasi?'
]

const scrollToBottom = async () => {
  await nextTick()
  if (chatBodyRef.value) {
    chatBodyRef.value.scrollTop = chatBodyRef.value.scrollHeight
  }
}

const sendQuickPrompt = (promptText) => {
  inputMessage.value = promptText
  sendMessage()
}

const sendMessage = async () => {
  const text = inputMessage.value.trim()
  if (!text || loading.value) return

  messages.value.push({ role: 'user', content: text })
  inputMessage.value = ''
  scrollToBottom()

  loading.value = true

  try {
    const res = await apiClient.post('/ai/chat', { message: text })
    const reply = res.data?.data?.reply || 'Maaf, terjadi kendala saat memproses jawaban AI.'
    messages.value.push({ role: 'ai', content: reply })
  } catch (err) {
    messages.value.push({ 
      role: 'ai', 
      content: 'Maaf, gagal menghubungkan ke server AI Assistant. Silakan coba beberapa saat lagi.' 
    })
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

const formatMessage = (content) => {
  if (!content) return ''
  let html = content
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/\n/g, '<br/>')
  return html
}

watch(isOpen, (newVal) => {
  if (newVal) {
    scrollToBottom()
  }
})
</script>

<style scoped>
.ai-widget-wrapper {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
}

.ai-fab-btn {
  position: relative;
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.75rem 1.25rem;
  border-radius: 9999px;
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-weight: 700;
  font-size: 0.9rem;
  box-shadow: 0 8px 25px rgba(99, 102, 241, 0.4);
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.ai-fab-btn:hover {
  transform: translateY(-4px) scale(1.03);
  box-shadow: 0 12px 30px rgba(99, 102, 241, 0.6);
}

.ai-fab-icon {
  font-size: 1.2rem;
}

.pulse-dot {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #34d399;
  box-shadow: 0 0 10px #34d399;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.7); }
  70% { transform: scale(1); box-shadow: 0 0 0 8px rgba(52, 211, 153, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(52, 211, 153, 0); }
}

.ai-chat-window {
  position: absolute;
  bottom: 70px;
  right: 0;
  width: 380px;
  max-width: calc(100vw - 32px);
  height: 520px;
  max-height: calc(100vh - 100px);
  background: #0f172a;
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: popUp 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes popUp {
  from { opacity: 0; transform: translateY(20px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.chat-header {
  padding: 1rem 1.25rem;
  background: rgba(30, 41, 59, 0.8);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.ai-avatar {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.3), rgba(79, 70, 229, 0.5));
  border: 1px solid rgba(99, 102, 241, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
}

.chat-title {
  font-size: 0.95rem;
  font-weight: 800;
  color: var(--text-main);
}

.chat-status {
  font-size: 0.72rem;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 0.35rem;
  margin-top: 0.1rem;
}

.green-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #34d399;
}

.close-chat-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.1rem;
  cursor: pointer;
  padding: 0.2rem 0.4rem;
  border-radius: 6px;
}

.close-chat-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.chat-body {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.chat-msg {
  display: flex;
  gap: 0.6rem;
  max-width: 85%;
}

.msg-ai {
  align-self: flex-start;
}

.msg-user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.msg-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: rgba(99, 102, 241, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  flex-shrink: 0;
}

.msg-bubble {
  padding: 0.75rem 0.95rem;
  border-radius: 14px;
  font-size: 0.82rem;
  line-height: 1.45;
  color: var(--text-main);
}

.msg-ai .msg-bubble {
  background: rgba(30, 41, 59, 0.9);
  border: 1px solid var(--border-color);
  border-top-left-radius: 2px;
}

.msg-user .msg-bubble {
  background: #4f46e5;
  color: white;
  border-top-right-radius: 2px;
}

.typing-bubble {
  display: flex;
  gap: 0.3rem;
  align-items: center;
  padding: 0.6rem 1rem;
}

.dot-typing {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-muted);
  animation: blink 1.4s infinite ease-in-out both;
}

.dot-typing:nth-child(1) { animation-delay: -0.32s; }
.dot-typing:nth-child(2) { animation-delay: -0.16s; }

@keyframes blink {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}

.chips-scroll {
  display: flex;
  gap: 0.4rem;
  padding: 0.5rem 0.75rem;
  overflow-x: auto;
  background: rgba(15, 23, 42, 0.5);
  border-top: 1px solid var(--border-color);
}

.chips-scroll::-webkit-scrollbar {
  height: 0px;
}

.chip-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  font-size: 0.72rem;
  padding: 0.35rem 0.65rem;
  border-radius: 9999px;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.15s ease;
}

.chip-btn:hover {
  background: rgba(99, 102, 241, 0.15);
  border-color: var(--accent-primary);
  color: #818cf8;
}

.chat-footer {
  padding: 0.75rem 1rem;
  background: rgba(30, 41, 59, 0.8);
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.chat-input {
  flex: 1;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: 9999px;
  padding: 0.55rem 1rem;
  color: var(--text-main);
  font-size: 0.82rem;
  outline: none;
}

.chat-input:focus {
  border-color: var(--accent-primary);
}

.chat-send-btn {
  background: var(--accent-primary);
  border: none;
  color: white;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.chat-send-btn:hover:not(:disabled) {
  transform: scale(1.1);
  background: #4338ca;
}

.chat-send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.mt-1 { margin-top: 0.25rem; }
</style>
