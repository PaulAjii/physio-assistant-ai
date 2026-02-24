<template>
  <div class="min-h-screen flex items-center justify-center p-6">
    <UCard class="w-full max-w-lg">
      <template #header>
        <div class="flex-flex-col gap-1">
          <h1 class="text-xl font-semibold">Physio Assistant</h1>
          <p class="text-sm text-gray-500">
            Upload a consultation audio recording to begin
          </p>
        </div>
      </template>

      <div class="flex flex-col gap-6">
        <UFormField label="Presenting Complaint" required>
          <USelect
            v-model="presentingComplaint"
            :items="complaints"
            class="w-full"
            placeholder="Select a presenting complaint"
          />
        </UFormField>

        <UFormField label="Consultation Audio" required>
          <div
            class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center cursor-pointer hover:border-primary transition-colors"
            @click="triggerFileInput"
            @dragover.prevent
            @drop.prevent="handleDrop"
          >
            <UIcon
              name="i-heroicons-musical-note"
              class="text-4xl text-gray-400 mx-auto mb-2"
            />
            <p class="text-sm text-gray-500">
              {{
                selectedFile
                  ? selectedFile.name
                  : "Click or drag and drop your audio file here"
              }}
            </p>
            <p class="text-xs text-gray-400 mt-1">
              MP3, WAV, M4A, WEBM supported
            </p>
            <input
              ref="fileInput"
              type="file"
              accept="audio/*"
              class="hidden"
              @change="handleFileChange"
            />
          </div>
        </UFormField>
      </div>

      <template #footer>
        <div class="flex justify-end">
          <UButton
            :loading="uploading"
            :disabled="!selectedFile || !presentingComplaint"
            label="Start Consultation"
            icon="i-heroicons-arrow-right"
            trailing
            @click="handleUpload"
          />
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const router = useRouter();
const { uploadAudio } = useConsultation();

const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
const uploading = ref(false);
const presentingComplaint = ref("");

const triggerFileInput = () => fileInput.value?.click();

const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (target.files?.[0]) selectedFile.value = target.files[0];
};

const handleDrop = (e: DragEvent) => {
  if (e.dataTransfer?.files?.[0]) selectedFile.value = e.dataTransfer.files[0];
};

const complaints = [
  "Knee Pain",
  "Shoulder Pain",
  "Lower Back Pain",
  "Neck Pain",
  "Hip Pain",
];

const handleUpload = async () => {
  if (!selectedFile.value) return;
  uploading.value = true;

  try {
    const response = await uploadAudio(selectedFile.value);
    router.push(
      `/consultations/${response?.jobID}?complaint=${encodeURIComponent(presentingComplaint.value)}`,
    );
  } catch (error) {
    console.error(error);
  } finally {
    uploading.value = false;
  }
};
</script>
