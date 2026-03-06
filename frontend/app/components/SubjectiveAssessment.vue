<template>
    <div class="flex flex-col gap-4">
          <!-- Processing State -->
          <UCard v-if="processing">
            <div class="flex flex-col items-center gap-4 py-8">
              <UIcon name="i-heroicons-arrow-path" class="text-4xl text-primary animate-spin" />
              <div class="text-center">
                <p class="font-medium">Processing Audio</p>
                <p class="text-sm text-gray-500">This may take 3-5 minutes depending on the length of the recording</p>
              </div>
            </div>
          </UCard>

          <!-- Error State -->
          <UAlert v-else-if="error" color="error" icon="i-heroicons-exclamation-circle" title="Processing Failed" :description="error" />

          <!-- Result -->
          <template v-else-if="result && editableSubjective">
            <UCard>
              <template #header>
                <div class="flex items-center justify-between">
                  <h2 class="font-semibold">Subjective Assessment</h2>
                  <UBadge color="warning" label="AI Draft" variant="subtle" />
                </div>

                <div class="flex items-center gap-2">
                  <UBadge :color="statusColor" :label="statusLabel" size="lg" />
                  <UButton
                    :icon="isEditing ? 'i-heroicons-check' : 'heroicons-pencil'"
                    :label="isEditing ? 'Done' : 'Edit'"
                    :color="isEditing ? 'success' : 'neutral'"
                    variant="soft"
                    size="sm"
                    @click="toggleEdit"
                  />
                </div>
              </template>

              <div class="flex flex-col gap-4">
                <!-- Presenting Complaint -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Presenting Complaint</h2>
                  </template>
                  <UTextarea
                    v-model="editableSubjective.presenting_complaint"
                    :disabled="!isEditing"
                    autoresize
                    class="w-full"
                  />
                </UCard>

                <!-- History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">History of Complaint</h2>
                  </template>
                  <UTextarea
                    v-model="editableSubjective.history_of_complaint"
                    :disabled="!isEditing"
                    autoresize
                    class="w-full"
                  />
                </UCard>

                <!-- Pain Profile -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Pain Profile</h2>
                  </template>
                  <div class="flex flex-col gap-3 text-sm">
                    <div class="flex justify-between items-center">
                      <span class="text-gray-500">Intensity</span>
                      <UInput
                        v-if="isEditing"
                        v-model="editableSubjective.pain_profile.intensity"
                        type="number"
                        :min="0"
                        :max="10"
                        size="sm"
                        class="w-20"
                      />
                      <UBadge v-else :color="painColor(editableSubjective.pain_profile.intensity)" :label="`${editableSubjective.pain_profile.intensity}/10`" />
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-gray-500">Quality</span>
                      <UInput
                        v-if="isEditing"
                        v-model="editableSubjective.pain_profile.quality"
                        size="sm"
                        class="w-48"
                      />
                      <span v-else>{{ editableSubjective.pain_profile.quality }}</span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-gray-500">Duration</span>
                      <UInput
                        v-if="isEditing"
                        v-model="editableSubjective.pain_profile.duration"
                        size="sm"
                        class="w-48"
                      />
                      <span v-else>{{ editableSubjective.pain_profile.duration }}</span>
                    </div>
                    <div class="flex flex-col gap-1">
                      <p class="text-gray-500">Location</p>
                      <UInput
                        v-if="isEditing"
                        :model-value="editableSubjective.pain_profile.location.join(', ')"
                        placeholder="Comma separated"
                        size="sm"
                        @update:model-value="editableSubjective.pain_profile.location = splitList($event)"
                      />
                      <div v-else class="flex flex-wrap gap-1">
                        <UBadge v-for="loc in editableSubjective.pain_profile.location" :key="loc" :label="loc" color="secondary" variant="subtle" />
                      </div>
                    </div>
                    <div class="flex flex-col gap-1">
                      <p class="text-gray-500">Aggravating Factors</p>
                      <UInput
                        v-if="isEditing"
                        :model-value="editableSubjective.pain_profile.aggravating.join(', ')"
                        placeholder="Comma separated"
                        size="sm"
                        @update:model-value="editableSubjective.pain_profile.aggravating = splitList($event)"
                      />
                      <div v-else class="flex flex-wrap gap-1">
                        <UBadge v-for="factor in editableSubjective.pain_profile.aggravating" :key="factor" :label="factor" color="error" variant="subtle" />
                      </div>
                    </div>
                    <div class="flex flex-col gap-1">
                      <p class="text-gray-500">Alleviating Factors</p>
                      <UInput
                        v-if="isEditing"
                        :model-value="editableSubjective.pain_profile.alleviating.join(', ')"
                        placeholder="Comma separated"
                        size="sm"
                        @update:model-value="editableSubjective.pain_profile.alleviating = splitList($event)"
                      />
                      <div v-else class="flex flex-wrap gap-1">
                        <UBadge v-for="factor in editableSubjective.pain_profile.alleviating" :key="factor" :label="factor" color="success" variant="subtle" />
                      </div>
                    </div>
                  </div>
                </UCard>

                <!-- Red Flags -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Red Flags</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="editableSubjective.red_flags.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    class="w-full"
                    @update:model-value="editableSubjective.red_flags = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="flag in editableSubjective.red_flags" :key="flag" class="text-sm flex items-center gap-2 w-full">
                        <UIcon name="i-heroicons-exclamation-triangle" class="text-error" />
                        {{ flag }}
                      </li>
                    </ul>
                    <p v-if="editableSubjective.red_flags.length === 0" class="text-sm text-gray-500">No red flags identified</p>
                  </template>
                </UCard>

                <!-- Associated Symptoms -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Associated Symptoms</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="editableSubjective.associated_symptoms.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="editableSubjective.associated_symptoms = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="symptom in editableSubjective.associated_symptoms" :key="symptom" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-check-circle" class="text-success" />
                        {{ symptom }}
                      </li>
                    </ul>
                    <p v-if="editableSubjective.associated_symptoms.length === 0" class="text-sm text-gray-500">No associated symptoms identified</p>
                  </template>
                </UCard>

                <!-- Relevant Medical History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Relevant Medical History</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="editableSubjective.relevant_medical_history.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="editableSubjective.relevant_medical_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="condition in editableSubjective.relevant_medical_history" :key="condition" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-check-circle" class="text-success" />
                        {{ condition }}
                      </li>
                    </ul>
                    <p v-if="editableSubjective.relevant_medical_history.length === 0" class="text-sm text-gray-500">No relevant medical history identified</p>
                  </template>
                </UCard>

                <!-- Past Surgical History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Past Surgical History</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="editableSubjective.past_surgical_history.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="editableSubjective.past_surgical_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="surgery in editableSubjective.past_surgical_history" :key="surgery" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-check-circle" class="text-success" />
                        {{ surgery }}
                      </li>
                    </ul>
                    <p v-if="editableSubjective.past_surgical_history.length === 0" class="text-sm text-gray-500">No past surgical history identified</p>
                  </template>
                </UCard>

                <!-- Drug History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Drug History</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="editableSubjective.drug_history.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="editableSubjective.drug_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="drug in editableSubjective.drug_history" :key="drug" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-check-circle" class="text-success" />
                        {{ drug }}
                      </li>
                    </ul>
                  </template>
                </UCard>

                <!-- Social History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Social History</h2>
                  </template>
                  <UTextarea
                    v-if="isEditing"
                    :model-value="editableSubjective.social_history.join(', ')"
                    placeholder="Comma separated"
                    :rows="3"
                    size="sm"
                    @update:model-value="editableSubjective.social_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="item in editableSubjective.social_history" :key="item" class="text-sm flex items-baseline gap-2">
                        <UIcon name="i-heroicons-user" class="text-gray-400 shrink-0" />
                        {{ item }}
                      </li>
                    </ul>
                  </template>
                </UCard>

                <!-- Family History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Family History</h2>
                  </template>
                  <UTextarea
                    v-if="isEditing"
                    v-model="editableSubjective.family_history"
                    :rows="2"
                    size="sm"
                    class="w-full"
                  />
                  <p v-else class="text-sm">{{ editableSubjective.family_history }}</p>
                </UCard>
              </div>
            </UCard>
          </template>
        </div>
</template>

<script setup lang="ts">
import type { AIData, Subjective } from "~~/types/assessment.d.ts"
const props = defineProps<{
  result: AIData | null;
  statusColor: "error" | "primary" | "secondary" | "success" | "info" | "warning" | "neutral" | undefined;
  statusLabel: string;
  processing: boolean;
  error: string | null;
}>()

const emit = defineEmits<{
  'update:subjective': [subjective: Subjective]
}>();

const isEditing = ref(false)

const editableSubjective = ref<Subjective | null>(
  props.result
    ? {
      ...props.result.subjective,
      pain_profile: {
        ...props.result.subjective.pain_profile,
      }
    }
    : null
);

watch(
  () => props.result,
  (newResult) => {
    if (newResult) {
      editableSubjective.value = {
        ...newResult.subjective,
        pain_profile: {
          ...newResult.subjective.pain_profile,
        }
      }
    }
  }
)

const splitList = (value: string) => value.split(',').map((s) => s.trim()).filter(Boolean)


const painColor = (intensity: number) => {
  if (intensity <= 3) return 'success'
  if (intensity <= 6) return 'warning'
  return 'error'
}

const toggleEdit = () => {
  isEditing.value = !isEditing.value;
  if (isEditing.value && editableSubjective.value) {
    emit('update:subjective', editableSubjective.value)
  }
}
</script>
