<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="toggleDropdown"
      :disabled="switching"
      class="flex items-center gap-2 rounded-lg px-2.5 py-1.5 text-sm font-semibold text-zinc-600 transition-colors hover:bg-zinc-100/60 dark:text-zinc-300 dark:hover:bg-zinc-800/60"
      :title="currentLocale?.name"
    >
      <Icon name="globe" size="sm" class="text-zinc-500 dark:text-zinc-400" />
      <span class="text-xs">{{ currentLocale?.code.toUpperCase() }}</span>
      <Icon
        name="chevronDown"
        size="xs"
        class="text-zinc-400 transition-transform duration-200"
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <transition name="dropdown">
      <div
        v-if="isOpen"
        class="absolute right-0 z-50 mt-1 w-32 overflow-hidden rounded-xl border border-zinc-200/80 bg-white/95 shadow-lg backdrop-blur-md dark:border-zinc-800 dark:bg-zinc-900/95"
      >
        <button
          v-for="locale in availableLocales"
          :key="locale.code"
          :disabled="switching"
          @click="selectLocale(locale.code)"
          class="flex w-full items-center gap-2.5 px-3.5 py-2 text-xs font-semibold text-zinc-700 transition-colors hover:bg-zinc-100/80 dark:text-zinc-200 dark:hover:bg-zinc-800/80"
          :class="{
            'bg-zinc-50 text-zinc-950 dark:bg-zinc-800/50 dark:text-white':
              locale.code === currentLocaleCode
          }"
        >
          <span class="text-[10px] font-bold text-zinc-400 dark:text-zinc-500 w-4">{{ locale.code.toUpperCase() }}</span>
          <span>{{ locale.name }}</span>
          <Icon v-if="locale.code === currentLocaleCode" name="check" size="sm" class="ml-auto text-zinc-900 dark:text-white" />
        </button>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { setLocale, availableLocales } from '@/i18n'

const { locale } = useI18n()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const switching = ref(false)

const currentLocaleCode = computed(() => locale.value)
const currentLocale = computed(() => availableLocales.find((l) => l.code === locale.value))

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

async function selectLocale(code: string) {
  if (switching.value || code === currentLocaleCode.value) {
    isOpen.value = false
    return
  }
  switching.value = true
  try {
    await setLocale(code)
    isOpen.value = false
  } finally {
    switching.value = false
  }
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}
</style>
