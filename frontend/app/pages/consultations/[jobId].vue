<template>
  <div class="min-h-screen p-6">
    <div class="max-w-6xl mx-auto flex flex-col gap-6">

      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold">Consultation Report</h1>
          <p class="text-sm text-gray-500">{{ complaint }}</p>
        </div>
        <UBadge :color="statusColor" :label="statusLabel" size="lg" />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

        <!-- LEFT: Subjective Report -->
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
            <!-- Presenting Complaint -->
            <UCard>
              <template #header>
                <h2 class="font-semibold">Presenting Complaint</h2>
              </template>
              <p class="text-sm">{{ result.data.subjective.presenting_complaint }}</p>
            </UCard>

            <!-- Pain Profile -->
            <UCard>
              <template #header>
                <h2 class="font-semibold">Pain Profile</h2>
              </template>
              <div class="flex flex-col gap-3 text-sm">
                <div class="flex justify-between">
                  <span class="text-gray-500">Intensity</span>
                  <UBadge :color="painColor(result.data.subjective.pain_profile.intensity)" :label="`${result.data.subjective.pain_profile.intensity}/10`" />
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Quality</span>
                  <span>{{ result.data.subjective.pain_profile.quality }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Duration</span>
                  <span>{{ result.data.subjective.pain_profile.duration }}</span>
                </div>
                <div>
                  <p class="text-gray-500 mb-1">Location</p>
                  <div class="flex flex-wrap gap-1">
                    <UBadge v-for="loc in result.data.subjective.pain_profile.location" :key="loc" :label="loc" color="secondary" variant="subtle" />
                  </div>
                </div>
                <div>
                  <p class="text-gray-500 mb-1">Aggravating Factors</p>
                  <div class="flex flex-wrap gap-1">
                    <UBadge v-for="factor in result.data.subjective.pain_profile.aggravating" :key="factor" :label="factor" color="error" variant="subtle" />
                  </div>
                </div>
                <div>
                  <p class="text-gray-500 mb-1">Alleviating Factors</p>
                  <div class="flex flex-wrap gap-1">
                    <UBadge v-for="factor in result.data.subjective.pain_profile.alleviating" :key="factor" :label="factor" color="success" variant="subtle" />
                  </div>
                </div>
              </div>
            </UCard>

            <!-- History -->
            <UCard>
              <template #header>
                <h2 class="font-semibold">History of Complaint</h2>
              </template>
              <p class="text-sm leading-relaxed">{{ result.data.subjective.history_of_complaint }}</p>
            </UCard>

            <!-- Drug History -->
            <UCard>
              <template #header>
                <h2 class="font-semibold">Drug History</h2>
              </template>
              <ul class="flex flex-col gap-1">
                <li v-for="drug in result.data.subjective.drug_history" :key="drug" class="text-sm flex items-center gap-2">
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
                <li v-for="item in result.data.subjective.social_history" :key="item" class="text-sm flex items-center gap-2">
                  <UIcon name="i-heroicons-user" class="text-gray-400" />
                  {{ item }}
                </li>
              </ul>
            </UCard>
          </template>
        </div>

        <!-- RIGHT: Objective Assessment -->
        <div class="flex flex-col gap-4">
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <h2 class="font-semibold">Objective Assessment</h2>
                <UBadge color="warning" label="In Progress" variant="subtle" />
              </div>
            </template>

            <div class="flex flex-col gap-4">
              <div v-for="(group, index) in objectiveChecklist" :key="index" class="flex flex-col gap-2">
                <div class="flex items-center gap-2">
                  <UBadge :color="priorityColor(group.priority)" :label="group.priority" variant="subtle" size="sm" />
                  <p class="font-medium text-sm">{{ group.category }}</p>
                </div>
                <div v-for="test in group.tests" :key="test" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <span class="text-sm">{{ test }}</span>
                  <div class="flex gap-2">
                    <UButton
                      size="xs"
                      :color="objectiveFindings[test] === 'positive' ? 'error' : 'neutral'"
                      variant="soft"
                      label="Positive"
                      @click="setFinding(test, 'positive')"
                    />
                    <UButton
                      size="xs"
                      :color="objectiveFindings[test] === 'negative' ? 'success' : 'neutral'"
                      variant="soft"
                      label="Negative"
                      @click="setFinding(test, 'negative')"
                    />
                    <UButton
                      size="xs"
                      :color="objectiveFindings[test] === 'not_tested' ? 'warning' : 'neutral'"
                      variant="soft"
                      label="N/T"
                      @click="setFinding(test, 'not_tested')"
                    />
                  </div>
                </div>

                <!-- Notes field -->
                <UTextarea
                  v-model="objectiveNotes[test]"
                  v-for="test in group.tests"
                  :key="`note-${test}`"
                  :placeholder="`Notes for ${test}...`"
                  size="sm"
                  class="mt-1"
                />
              </div>
            </div>

            <template #footer>
              <div class="flex justify-end">
                <UButton
                  label="Submit Objective Assessment"
                  icon="i-heroicons-check"
                  :disabled="!result || objectiveEmpty"
                  @click="submitObjective"
                />
              </div>
            </template>
          </UCard>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { streamResult } = useConsultation()

const jobId = route.params.jobId as string
const complaint = decodeURIComponent(route.query.complaint as string || '')

const processing = ref(true)
const error = ref<string | null>(null)
const result = ref<any>(null)

const objectiveFindings = ref<Record<string, string>>({})
const objectiveNotes = ref<Record<string, string>>({})

// Preset checklists per complaint
const presets: Record<string, { category: string; tests: string[]; priority: string }[]> = {
  'Knee Pain': [
    { category: 'Gait Analysis', tests: ['Observation of gait pattern', 'Antalgic gait check'], priority: 'high' },
    { category: 'Range of Motion', tests: ['Active knee flexion', 'Active knee extension', 'Passive knee flexion', 'Passive knee extension'], priority: 'high' },
    { category: 'Strength Testing', tests: ['Quadriceps', 'Hamstrings'], priority: 'high' },
    { category: 'Special Tests', tests: ["McMurray's", "Lachman's", 'Anterior Drawer', 'Posterior Drawer', 'Patellar grind'], priority: 'high' },
    { category: 'Palpation', tests: ['Joint line tenderness', 'Patellar tenderness', 'Surrounding musculature'], priority: 'medium' },
    { category: 'Neurological Screen', tests: ['Dermatomes', 'Myotomes', 'Reflexes'], priority: 'medium' },
  ],
  'Shoulder Pain': [
    { category: 'Range of Motion', tests: ['Flexion', 'Abduction', 'Internal rotation', 'External rotation'], priority: 'high' },
    { category: 'Special Tests', tests: ["Hawkins-Kennedy", "Neer's", "Empty Can", "Speed's"], priority: 'high' },
    { category: 'Strength Testing', tests: ['Rotator cuff', 'Deltoid', 'Biceps'], priority: 'high' },
  ],
  'Lower Back Pain': [
    { category: 'Posture & Gait', tests: ['Postural assessment', 'Gait observation'], priority: 'high' },
    { category: 'Range of Motion', tests: ['Flexion', 'Extension', 'Lateral flexion', 'Rotation'], priority: 'high' },
    { category: 'Special Tests', tests: ["Straight Leg Raise", "FABER", "FADIR", "Slump test"], priority: 'high' },
    { category: 'Neurological Screen', tests: ['Dermatomes L1-S1', 'Myotomes', 'Reflexes'], priority: 'high' },
  ],
  'Neck Pain': [
    { category: 'Range of Motion', tests: ['Flexion', 'Extension', 'Rotation left', 'Rotation right', 'Lateral flexion'], priority: 'high' },
    { category: 'Special Tests', tests: ["Spurling's", "Distraction test", "Upper limb tension test"], priority: 'high' },
    { category: 'Neurological Screen', tests: ['Dermatomes C4-T1', 'Myotomes', 'Reflexes'], priority: 'high' },
  ],
}

const objectiveChecklist = computed(() => presets[complaint] || [])

const objectiveEmpty = computed(() => Object.keys(objectiveFindings.value).length === 0)

const statusColor = computed(() => {
  if (processing.value) return 'warning'
  if (error.value) return 'error'
  return 'success'
})

const statusLabel = computed(() => {
  if (processing.value) return 'Processing'
  if (error.value) return 'Failed'
  return 'Complete'
})

const painColor = (intensity: number) => {
  if (intensity <= 3) return 'success'
  if (intensity <= 6) return 'warning'
  return 'error'
}

const priorityColor = (priority: string) => {
  if (priority === 'high') return 'error'
  if (priority === 'medium') return 'warning'
  return 'success'
}

const setFinding = (test: string, value: string) => {
  objectiveFindings.value[test] = value
}

const submitObjective = () => {
  const findings = objectiveChecklist.value.flatMap(group =>
    group.tests.map(test => ({
      category: group.category,
      test,
      result: objectiveFindings.value[test] || 'not_tested',
      notes: objectiveNotes.value[test] || ''
    }))
  )
  console.log('Objective findings:', findings)
  // TODO: send to Go backend
}

// Start SSE stream on mount
onMounted(() => {
  streamResult(
    jobId,
    (data) => {
      result.value = data
      processing.value = false
    },
    (err) => {
      error.value = err.message || 'An error occurred during processing'
      processing.value = false
    }
  )
})
</script>
