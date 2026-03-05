<template>
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
                <div v-for="test in group.tests" :key="test" class="flex flex-col gap-2 p-3 bg-gray-700 rounded-lg">
                  <div class="flex items-center justify-between">
                    <span class="text-sm">{{ test }}</span>
                    <div class="flex gap-2">
                      <UButton
                        size="xs"
                        :color="objectiveFindings[test]?.result === 'positive' ? 'error' : 'neutral'"
                        variant="soft"
                        label="Positive"
                        @click="setFinding(test, group.category, 'positive')"
                      />
                      <UButton
                        size="xs"
                        :color="objectiveFindings[test]?.result === 'negative' ? 'success' : 'neutral'"
                        variant="soft"
                        label="Negative"
                        @click="setFinding(test, group.category, 'negative')"
                      />
                      <UButton
                        size="xs"
                        :color="objectiveFindings[test]?.result === 'not_tested' ? 'warning' : 'neutral'"
                        variant="soft"
                        label="N/T"
                        @click="setFinding(test, group.category, 'not_tested')"
                      />
                    </div>
                  </div>
                  <UTextarea
                    :model-value="objectiveFindings[test]?.notes ?? ''"
                    :placeholder="`Notes for ${test}...`"
                    size="sm"
                    @update:model-value="updateNotes(test, group.category, $event)"
                  />
                </div>
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
</template>

<script setup lang="ts">
import type { ObjectiveFinding } from '~~/types/assessment';

const props = defineProps<{
    result: any;
    complaint: string;
}>();

const objectiveFindings = ref<Record<string, ObjectiveFinding>>({})

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

const objectiveChecklist = computed(() => presets[props.complaint] || [])
const objectiveEmpty = computed(() => Object.keys(objectiveFindings.value).length === 0)

const priorityColor = (priority: string) => {
  if (priority === 'high') return 'error'
  if (priority === 'medium') return 'warning'
  return 'success'
}

const setFinding = (test: string, value: string, category: string) => {
    if (!objectiveFindings.value[test]) {
      objectiveFindings.value[test] = { category, test, result: value, notes: ""}
    } else {
      objectiveFindings.value[test].result = value;
    }
}

const emit = defineEmits<{
  'submit:objective': [findings: ObjectiveFinding[]];
}>();

const updateNotes = (test: string, category: string, value: string) => {
  if (!objectiveFindings.value[test]) {
    objectiveFindings.value[test] = { category, test, result: 'not_tested', notes: value }
  } else {
    objectiveFindings.value[test].notes = value
  }
}

const submitObjective = () => {
  const findings = Object.values(objectiveFindings.value);
  emit('submit:objective', findings);
}
</script>
