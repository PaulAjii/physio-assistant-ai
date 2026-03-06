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
                <div v-for="test in group.tests" :key="test.name" class="flex flex-col gap-2 p-3 bg-gray-700 rounded-lg">
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-200">{{ test.name }}</span>
                    <div v-if="test.test === 'binary'" class="flex gap-2">
                      <UButton
                        size="xs"
                        :color="getResult(test.name) === 'positive' ? 'error' : 'neutral'"
                        variant="soft"
                        label="Positive"
                        @click="setFinding(test, group.category, 'positive')"
                      />
                      <UButton
                        size="xs"
                        :color="getResult(test.name) === 'negative' ? 'success' : 'neutral'"
                        variant="soft"
                        label="Negative"
                        @click="setFinding(test, group.category, 'negative')"
                      />
                      <UButton
                        size="xs"
                        :color="getResult(test.name) === 'not_tested' ? 'warning' : 'neutral'"
                        variant="soft"
                        label="N/T"
                        @click="setFinding(test, group.category, 'not_tested')"
                      />
                    </div>

                    <div v-else-if="test.test === 'measurement'" class="flex items-center gap-2">
                      <UInput
                        :model-value="getValue(test.name)"
                        type="number"
                        size="sm"
                        class="w-24"
                        :min="0"
                        :max="360"
                        :placeholder="test.unit ?? ''"
                        @update:model-value="updateValue(test, group.category, $event)"
                      />
                      <span v-if="test.unit" class="text-xs text-gray-400">{{ test.unit  }}</span>
                    </div>

                    <div v-else-if="test.test === 'grading'" class="flex gap-2">
                      <UButton
                        v-for="grade in [0, 1, 2, 3, 4 , 5]"
                        :key="grade"
                        size="xs"
                        :color="getValue(test.name) === String(grade) ? 'primary' : 'neutral'"
                        variant="soft"
                        :label="String(grade)"
                        @click="updateValue(test, group.category, String(grade))"
                      />
                    </div>
                  </div>
                  <UTextarea
                    :model-value="getNotes(test.name)"
                    class="field-sizing-content"
                    :placeholder="test.test === 'notes' ? `Findings for ${ test.name }` : `Additional notes for ${ test.name }`"
                    size="sm"
                    :rows="3"
                    autoresize
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
import type { ObjectiveFinding, ObjectiveTest } from '~~/types/assessment';
import { getObjectiveTemplateByComplaint } from '~/composables/templates';

const props = defineProps<{
    result: any;
    complaint: string;
}>();

const objectiveFindings = ref<Record<string, ObjectiveFinding>>({})

const objectiveChecklist = computed(() => {
  return getObjectiveTemplateByComplaint(props.complaint)?.categories ?? [];
})
const objectiveEmpty = computed(() => Object.keys(objectiveFindings.value).length === 0)

const priorityColor = (priority: string) => {
  if (priority === 'high') return 'error'
  if (priority === 'medium') return 'warning'
  return 'success'
}

const getResult = (name: string) => objectiveFindings.value[name]?.result ?? ''
const getValue = (name: string) => objectiveFindings.value[name]?.value ?? ''
const getNotes = (name: string) => objectiveFindings.value[name]?.notes ?? ''

const setFinding = (test: ObjectiveTest, category: string, result: string) => {
    if (!objectiveFindings.value[test.name]) {
      objectiveFindings.value[test.name] = { category, type: test.test, test: test.name, result, value: '', unit: test.unit ?? '', notes: ""}
    } else {
      objectiveFindings.value[test.name].result = result;
    }
}

const emit = defineEmits<{
  'submit:objective': [findings: ObjectiveFinding[]];
}>();

const updateValue = (test: ObjectiveTest, category: string, value: string) => {
  if (!objectiveFindings.value[test.name]) {
    objectiveFindings.value[test.name] = { category, type: test.test, test: test.name, result: '', value, unit: test.unit ?? '', notes: '' }
  } else {
    objectiveFindings.value[test.name].value = value;
  }
}

const updateNotes = (test: ObjectiveTest, category: string, notes: string) => {
  if (!objectiveFindings.value[test.name]) {
    objectiveFindings.value[test.name] = { category, test: test.name, result: '', type: test.test, value: '', unit: test.unit ?? '', notes }
  } else {
    objectiveFindings.value[test.name].notes = notes
  }
}

const submitObjective = () => {
  const findings = Object.values(objectiveFindings.value);
  emit('submit:objective', findings);
}
</script>
