<template>
  <div class="min-h-screen flex items-center justify-center p-6">
    <UCard class="w-full max-w-lg">
      <template #header>
        <div class="flex-flex-col gap-1">
          <h1 class="text-xl font-semibold">Physio Assistant</h1>
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup>
const router = useRouter();
const { uploadAudio } = useConsultation();
</script>
