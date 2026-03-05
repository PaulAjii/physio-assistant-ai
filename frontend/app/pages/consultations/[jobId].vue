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
        <SubjectiveAssessment
          :result="result"
          :statusColor="statusColor"
          :error="error"
          :statusLabel="statusLabel"
          :processing="processing"
          @update:subjective="correctedSubjective = $event"
        />

        <!-- RIGHT: Objective Assessment -->
        <ObjectiveAssessment
        :result="result"
        :complaint="complaint"
        @submit:objective="submitObjective" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AIData, ObjectiveFinding, Subjective } from '~~/types/assessment';

const route = useRoute();
const router = useRouter();
const { streamResult, submitAssessment } = useConsultation();

const jobId = route.params.jobId as string
const complaint = decodeURIComponent(route.query.complaint as string || '')
const result = ref<AIData | null>(null)
const processing = ref(true)
const error = ref<string | null>(null)

const correctedSubjective = ref<Subjective | null>(null);

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

const submitObjective = async (findings: ObjectiveFinding[]) => {
  if (!correctedSubjective.value) return

  try {
    const response = await submitAssessment({
      jobID: jobId,
      complaint,
      corrected_subjective: correctedSubjective.value,
      objective_findings: findings,
    });

    router.push(`/assessments/${ response.data.id }`)
  } catch (err: any) {
    console.error("Submission failed: ", err)
  }
}

// Start SSE stream on mount
onMounted(() => {
  streamResult(
    jobId,
    (data) => {
      result.value = data?.processedData
      correctedSubjective.value = data?.processedData?.subjective
      processing.value = false
    },
    (err) => {
      error.value = err.message || 'An error occurred during processing'
      processing.value = false
    }
  )
})
</script>
