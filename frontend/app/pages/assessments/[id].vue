<template>
  <div class="min-h-screen p-6">
    <div class="max-w-6xl mx-auto flex flex-col gap-6">

      <!-- Header -->
      <div class="flex items-center justify-between">
        <div class="flex flex-col gap-1">
          <h1 class="text-xl font-semibold">Assessment Report</h1>
          <p class="text-sm text-gray-500">{{ assessment?.complaint }}</p>
        </div>
        <div class="flex items-center gap-3">
          <UBadge color="success" label="Final" size="lg" />
          <UButton
            :icon="isEditing ? 'i-heroicons-check' : 'i-heroicons-pencil'"
            :label="isEditing ? 'Save Changes' : 'Edit'"
            :color="isEditing ? 'success' : 'neutral'"
            :loading="saving"
            variant="soft"
            @click="toggleEdit"
          />
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex flex-col items-center gap-4 py-16">
        <UIcon name="i-heroicons-arrow-path" class="text-4xl text-primary animate-spin" />
        <p class="text-sm text-gray-500">Loading assessment...</p>
      </div>

      <!-- Error State -->
      <UAlert
        v-else-if="fetchError"
        color="error"
        icon="i-heroicons-exclamation-circle"
        title="Failed to load assessment"
        :description="fetchError"
      />

      <!-- Assessment Content -->
      <template v-else-if="assessment">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">

          <!-- LEFT: Subjective -->
          <div class="flex flex-col gap-4">
            <UCard>
              <template #header>
                <h2 class="font-semibold">Subjective Assessment</h2>
              </template>

              <div class="flex flex-col gap-4">

                <!-- Presenting Complaint -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Presenting Complaint</h2>
                  </template>
                  <UTextarea
                    v-model="assessment.assessment.presenting_complaint"
                    :disabled="!isEditing"
                    :rows="2"
                    class="w-full"
                  />
                </UCard>

                <!-- History of Complaint -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">History of Complaint</h2>
                  </template>
                  <UTextarea
                    v-model="assessment.assessment.history_of_complaint"
                    :disabled="!isEditing"
                    :rows="5"
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
                        v-model.number="assessment.assessment.pain_profile.intensity"
                        type="number"
                        :min="0"
                        :max="10"
                        size="sm"
                        class="w-20"
                      />
                      <UBadge
                        v-else
                        :color="painColor(assessment.assessment.pain_profile.intensity)"
                        :label="`${assessment.assessment.pain_profile.intensity}/10`"
                      />
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-gray-500">Quality</span>
                      <UInput
                        v-if="isEditing"
                        v-model="assessment.assessment.pain_profile.quality"
                        size="sm"
                        class="w-48"
                      />
                      <span v-else>{{ assessment.assessment.pain_profile.quality }}</span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-gray-500">Duration</span>
                      <UInput
                        v-if="isEditing"
                        v-model="assessment.assessment.pain_profile.duration"
                        size="sm"
                        class="w-48"
                      />
                      <span v-else>{{ assessment.assessment.pain_profile.duration }}</span>
                    </div>
                    <div class="flex flex-col gap-1">
                      <p class="text-gray-500">Location</p>
                      <UInput
                        v-if="isEditing"
                        :model-value="assessment.assessment.pain_profile.location.join(', ')"
                        placeholder="Comma separated"
                        size="sm"
                        @update:model-value="assessment.assessment.pain_profile.location = splitList($event)"
                      />
                      <div v-else class="flex flex-wrap gap-1">
                        <UBadge v-for="loc in assessment.assessment.pain_profile.location" :key="loc" :label="loc" color="secondary" variant="subtle" />
                      </div>
                    </div>
                    <div class="flex flex-col gap-1">
                      <p class="text-gray-500">Aggravating Factors</p>
                      <UInput
                        v-if="isEditing"
                        :model-value="assessment.assessment.pain_profile.aggravating.join(', ')"
                        placeholder="Comma separated"
                        size="sm"
                        @update:model-value="assessment.assessment.pain_profile.aggravating = splitList($event)"
                      />
                      <div v-else class="flex flex-wrap gap-1">
                        <UBadge v-for="factor in assessment.assessment.pain_profile.aggravating" :key="factor" :label="factor" color="error" variant="subtle" />
                      </div>
                    </div>
                    <div class="flex flex-col gap-1">
                      <p class="text-gray-500">Alleviating Factors</p>
                      <UInput
                        v-if="isEditing"
                        :model-value="assessment.assessment.pain_profile.alleviating.join(', ')"
                        placeholder="Comma separated"
                        size="sm"
                        @update:model-value="assessment.assessment.pain_profile.alleviating = splitList($event)"
                      />
                      <div v-else class="flex flex-wrap gap-1">
                        <UBadge v-for="factor in assessment.assessment.pain_profile.alleviating" :key="factor" :label="factor" color="success" variant="subtle" />
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
                    :model-value="assessment.assessment.red_flags.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="assessment.assessment.red_flags = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="flag in assessment.assessment.red_flags" :key="flag" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-exclamation-triangle" class="text-error" />
                        {{ flag }}
                      </li>
                    </ul>
                    <p v-if="assessment.assessment.red_flags.length === 0" class="text-sm text-gray-500">No red flags identified</p>
                  </template>
                </UCard>

                <!-- Associated Symptoms -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Associated Symptoms</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="assessment.assessment.associated_symptoms.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="assessment.assessment.associated_symptoms = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="symptom in assessment.assessment.associated_symptoms" :key="symptom" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-check-circle" class="text-success" />
                        {{ symptom }}
                      </li>
                    </ul>
                    <p v-if="assessment.assessment.associated_symptoms.length === 0" class="text-sm text-gray-500">No associated symptoms identified</p>
                  </template>
                </UCard>

                <!-- Relevant Medical History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Relevant Medical History</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="assessment.assessment.relevant_medical_history.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="assessment.assessment.relevant_medical_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="condition in assessment.assessment.relevant_medical_history" :key="condition" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-check-circle" class="text-success" />
                        {{ condition }}
                      </li>
                    </ul>
                    <p v-if="assessment.assessment.relevant_medical_history.length === 0" class="text-sm text-gray-500">No relevant medical history identified</p>
                  </template>
                </UCard>

                <!-- Past Surgical History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Past Surgical History</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="assessment.assessment.past_surgical_history.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="assessment.assessment.past_surgical_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="surgery in assessment.assessment.past_surgical_history" :key="surgery" class="text-sm flex items-center gap-2">
                        <UIcon name="i-heroicons-check-circle" class="text-success" />
                        {{ surgery }}
                      </li>
                    </ul>
                    <p v-if="assessment.assessment.past_surgical_history.length === 0" class="text-sm text-gray-500">No past surgical history identified</p>
                  </template>
                </UCard>

                <!-- Drug History -->
                <UCard>
                  <template #header>
                    <h2 class="font-semibold">Drug History</h2>
                  </template>
                  <UInput
                    v-if="isEditing"
                    :model-value="assessment.assessment.drug_history.join(', ')"
                    placeholder="Comma separated"
                    size="sm"
                    @update:model-value="assessment.assessment.drug_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="drug in assessment.assessment.drug_history" :key="drug" class="text-sm flex items-center gap-2">
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
                    :model-value="assessment.assessment.social_history.join(', ')"
                    placeholder="Comma separated"
                    :rows="3"
                    size="sm"
                    @update:model-value="assessment.assessment.social_history = splitList($event)"
                  />
                  <template v-else>
                    <ul class="flex flex-col gap-1">
                      <li v-for="item in assessment.assessment.social_history" :key="item" class="text-sm flex items-baseline gap-2">
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
                    v-model="assessment.assessment.family_history"
                    :rows="2"
                    size="sm"
                    class="w-full"
                  />
                  <p v-else class="text-sm">{{ assessment.assessment.family_history }}</p>
                </UCard>

              </div>
            </UCard>

            <!-- AI Draft Reference -->
            <UCard>
              <template #header>
                <div class="flex items-center gap-2">
                  <h2 class="font-semibold">AI Draft Reference</h2>
                  <UBadge color="warning" label="Read Only" variant="subtle" />
                </div>
              </template>
              <div class="flex flex-col gap-3">
                <div class="flex flex-col gap-1">
                  <label class="text-xs text-gray-400">Original AI Presenting Complaint</label>
                  <p class="text-sm text-gray-500 italic leading-relaxed">
                    {{ assessment.ai_draft.data.subjective.presenting_complaint }}
                  </p>
                </div>
                <div class="flex items-center gap-4 text-xs text-gray-400">
                  <span>Language: {{ assessment.ai_draft.data.meta.language_detected }}</span>
                  <UDivider orientation="vertical" />
                  <span>Speaker ID: {{ assessment.ai_draft.data.meta.speaker_identification }}</span>
                </div>
              </div>
            </UCard>
          </div>

          <!-- RIGHT: Objective -->
          <div class="flex flex-col gap-4">
            <UCard>
              <template #header>
                <h2 class="font-semibold">Objective Findings</h2>
              </template>

              <div class="flex flex-col gap-3">
                <div
                  v-for="finding in assessment.assessment.objective_findings"
                  :key="finding.test"
                  class="flex flex-col gap-2 p-3 bg-gray-700 rounded-lg"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex flex-col gap-0.5">
                      <span class="text-sm font-medium">{{ finding.test }}</span>
                      <span class="text-xs text-gray-400">{{ finding.category }}</span>
                    </div>
                    <div v-if="!isEditing">
                      <UBadge
                        :color="findingColor(finding.result)"
                        :label="finding.result.replace('_', ' ')"
                        variant="subtle"
                        size="sm"
                      />
                    </div>
                    <div v-else class="flex gap-1">
                      <UButton
                        size="xs"
                        :color="finding.result === 'positive' ? 'error' : 'neutral'"
                        variant="soft"
                        label="Positive"
                        @click="finding.result = 'positive'"
                      />
                      <UButton
                        size="xs"
                        :color="finding.result === 'negative' ? 'success' : 'neutral'"
                        variant="soft"
                        label="Negative"
                        @click="finding.result = 'negative'"
                      />
                      <UButton
                        size="xs"
                        :color="finding.result === 'not_tested' ? 'warning' : 'neutral'"
                        variant="soft"
                        label="N/T"
                        @click="finding.result = 'not_tested'"
                      />
                    </div>
                  </div>
                  <UInput
                    v-model="finding.notes"
                    :disabled="!isEditing"
                    placeholder="Notes..."
                    size="sm"
                  />
                </div>
              </div>
            </UCard>

            <!-- Assessment Info -->
            <UCard>
              <template #header>
                <h2 class="font-semibold">Assessment Info</h2>
              </template>
              <div class="flex flex-col gap-2 text-sm">
                <div class="flex justify-between">
                  <span class="text-gray-500">Assessment ID</span>
                  <span class="font-mono text-xs">{{ assessment.id }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Created</span>
                  <span>{{ formatDate(assessment.created_at) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Complaint</span>
                  <span>{{ assessment.complaint }}</span>
                </div>
              </div>
            </UCard>
          </div>

        </div>
      </template>

    </div>
  </div>
</template>

<script setup lang="ts">
import type { CollatedAssessment } from '~~/types/assessment'

const route = useRoute()
const { getAssessment, updateAssessment } = useConsultation()

const id = route.params.id as string
const CACHE_KEY = `assessment:${id}`

const loading = ref(true)
const fetchError = ref<string | null>(null)
const assessment = ref<CollatedAssessment | null>(null)
const isEditing = ref(false)
const saving = ref(false)

const splitList = (value: string) =>
  value.split(',').map((s) => s.trim()).filter(Boolean)

const formatDate = (dateStr: string) =>
  new Date(dateStr).toLocaleString()

const painColor = (intensity: number) => {
  if (intensity <= 3) return 'success'
  if (intensity <= 6) return 'warning'
  return 'error'
}

const findingColor = (result: string) => {
  if (result === 'positive') return 'error'
  if (result === 'negative') return 'success'
  return 'warning'
}

const loadAssessment = async () => {
  loading.value = true
  fetchError.value = null

  try {
    // Check localStorage cache first
    const cached = localStorage.getItem(CACHE_KEY)
    if (cached) {
      assessment.value = JSON.parse(cached)
      loading.value = false
      return
    }

    // No cache — fetch from API
    const response = await getAssessment(id)
    assessment.value = response.data

    // Cache for future reloads
    localStorage.setItem(CACHE_KEY, JSON.stringify(response.data))
  } catch (error: any) {
    fetchError.value = error?.message || 'Failed to load assessment'
  } finally {
    loading.value = false
  }
}

const toggleEdit = async () => {
  if (!isEditing.value) {
    isEditing.value = true
    return
  }

  // Save changes
  saving.value = true
  try {
    const response = await updateAssessment(id, assessment.value!.assessment)
    assessment.value = response.data

    // Update cache with latest saved data
    localStorage.setItem(CACHE_KEY, JSON.stringify(response.data))

    isEditing.value = false
  } catch (error: any) {
    console.error('Failed to save:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadAssessment()
})
</script>
