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
          <template v-else-if="result">
            <UCard>
              <template #header>
                <div class="flex items-center justify-between">
                  <h2 class="font-semibold">Subjective Assessment</h2>
                  <UBadge :color="statusColor" :label="statusLabel" size="lg" />
                </div>
              </template>
            
              <div class="flex flex-col gap-4">
                <!-- Presenting Complaint -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Presenting Complaint</h2>
                  </template>
                  <p class="text-sm">{{ result.processedData.subjective.presenting_complaint }}</p>
                </UCard>

                <!-- History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">History of Complaint</h2>
                  </template>
                  <p class="text-sm leading-relaxed">{{ result.processedData.subjective.history_of_complaint }}</p>
                </UCard>

                <!-- Pain Profile -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Pain Profile</h2>
                  </template>
                  <div class="flex flex-col gap-3 text-sm">
                    <div class="flex justify-between">
                      <span class="text-gray-500">Intensity</span>
                      <UBadge :color="painColor(result.processedData.subjective.pain_profile.intensity)" :label="`${result.processedData.subjective.pain_profile.intensity}/10`" />
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">Quality</span>
                      <span>{{ result.processedData.subjective.pain_profile.quality }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">Duration</span>
                      <span>{{ result.processedData.subjective.pain_profile.duration }}</span>
                    </div>
                    <div>
                      <p class="text-gray-500 mb-1">Location</p>
                      <div class="flex flex-wrap gap-1">
                        <UBadge v-for="loc in result.processedData.subjective.pain_profile.location" :key="loc" :label="loc" color="secondary" variant="subtle" />
                      </div>
                    </div>
                    <div>
                      <p class="text-gray-500 mb-1">Aggravating Factors</p>
                      <div class="flex flex-wrap gap-1">
                        <UBadge v-for="factor in result.processedData.subjective.pain_profile.aggravating" :key="factor" :label="factor" color="error" variant="subtle" />
                      </div>
                    </div>
                    <div>
                      <p class="text-gray-500 mb-1">Alleviating Factors</p>
                      <div class="flex flex-wrap gap-1">
                        <UBadge v-for="factor in result.processedData.subjective.pain_profile.alleviating" :key="factor" :label="factor" color="success" variant="subtle" />
                      </div>
                    </div>
                  </div>
                </UCard>

                <!-- Red Flags -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Red Flags</h2>
                  </template>
                  <ul class="flex flex-col gap-1">
                    <li v-for="flag in result.processedData.subjective.red_flags" :key="flag" class="text-sm flex items-center gap-2">
                      <UIcon name="i-heroicons-check-circle" class="text-green-500" />
                      {{ flag }}
                    </li>
                  </ul>
                  <p v-if="result.processedData.subjective.red_flags.length === 0" class="text-sm text-gray-500">No red flags identified</p>
                </UCard>

                <!-- Associated Symptoms -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Associated Symptoms</h2>
                  </template>
                  <ul class="flex flex-col gap-1">
                    <li v-for="symptom in result.processedData.subjective.associated_symptoms" :key="symptom" class="text-sm flex items-center gap-2">
                      <UIcon name="i-heroicons-check-circle" class="text-green-500" />
                      {{ symptom }}
                    </li>
                  </ul>
                  <p v-if="result.processedData.subjective.associated_symptoms.length === 0" class="text-sm text-gray-500">No associated symptoms identified</p>
                </UCard>

                <!-- Relevant Medical History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Relevant Medical History</h2>
                  </template>
                  <ul class="flex flex-col gap-1">
                    <li v-for="condition in result.processedData.subjective.relevant_medical_history" :key="condition" class="text-sm flex items-center gap-2">
                      <UIcon name="i-heroicons-check-circle" class="text-green-500" />
                      {{ condition }}
                    </li>
                  </ul>
                  <p v-if="result.processedData.subjective.relevant_medical_history.length === 0" class="text-sm text-gray-500">No relevant medical history identified</p>
                </UCard>

                <!-- Past Surgical History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Past Surgical History</h2>
                  </template>
                  <ul class="flex flex-col gap-1">
                    <li v-for="surgery in result.processedData.subjective.past_surgical_history" :key="surgery" class="text-sm flex items-center gap-2">
                      <UIcon name="i-heroicons-check-circle" class="text-green-500" />
                      {{ surgery }}
                    </li>
                  </ul>
                  <p v-if="result.processedData.subjective.past_surgical_history.length === 0" class="text-sm text-gray-500">No past surgical history identified</p>
                </UCard>

                <!-- Drug History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Drug History</h2>
                  </template>
                  <ul class="flex flex-col gap-1">
                    <li v-for="drug in result.processedData.subjective.drug_history" :key="drug" class="text-sm flex items-center gap-2">
                      <UIcon name="i-heroicons-check-circle" class="text-green-500" />
                      {{ drug }}
                    </li>
                  </ul>
                </UCard>

                <!-- Social History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Social History</h2>
                  </template>
                  <ul class="flex flex-col gap-1">
                    <li v-for="item in result.processedData.subjective.social_history" :key="item" class="text-sm flex items-baseline gap-2">
                      <UIcon name="i-heroicons-check-circle" class="text-green-500" />
                      {{ item }}
                    </li>
                  </ul>
                </UCard>

                <!-- Family History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Family History</h2>
                  </template>
                  <p class="text-sm">{{ result.processedData.subjective.family_history }}</p>
                </UCard>
              </div>
            </UCard>
          </template>
        </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  result: any;
  statusColor: "error" | "primary" | "secondary" | "success" | "info" | "warning" | "neutral" | undefined;
  statusLabel: string;
  processing: boolean;
  error: string | null;
}>()

const painColor = (intensity: number) => {
  if (intensity <= 3) return 'success'
  if (intensity <= 6) return 'warning'
  return 'error'
}
</script>