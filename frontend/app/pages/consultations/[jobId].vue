<template>
  <div class="min-h-screen p-6">
    <div class="max-w-6xl mx-auto flex flex-col gap-6">

      <!-- Header -->
      <div>
          <h1 class="text-xl font-semibold">Clerking</h1>
          <p class="text-sm text-gray-500">{{ complaint }}</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

        <!-- LEFT: Subjective Report -->
        <SubjectiveAssessment :result="result" :statusColor="statusColor" :error="error" :statusLabel="statusLabel" :processing="processing" />

        <!-- RIGHT: Objective Assessment -->
        <ObjectiveAssessment :result="result" :complaint="complaint" @submit-objective="submitObjective" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { streamResult } = useConsultation()

const jobId = route.params.jobId as string
const complaint = decodeURIComponent(route.query.complaint as string || '')
const result = ref<any>(null)
const processing = ref(true)
const error = ref<string | null>(null)

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

const submitObjective = (findings: any[]) => {
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
