import type { ObjectiveTemplate } from '~~/types/assessment'

export const defaultObjectiveTemplates: ObjectiveTemplate[] = [
  {
    name: 'Knee Pain Assessment',
    complaint: 'Knee Pain',
    categories: [
      {
        category: 'Gait Analysis',
        priority: 'high',
        tests: [
          { name: 'Observation of gait pattern', test: 'notes' },
          { name: 'Antalgic gait check', test: 'binary' },
        ],
      },
      {
        category: 'Range of Motion',
        priority: 'high',
        tests: [
          { name: 'Active knee flexion', test: 'measurement', unit: 'degrees' },
          { name: 'Active knee extension', test: 'measurement', unit: 'degrees' },
          { name: 'Passive knee flexion', test: 'measurement', unit: 'degrees' },
          { name: 'Passive knee extension', test: 'measurement', unit: 'degrees' },
        ],
      },
      {
        category: 'Strength Testing',
        priority: 'high',
        tests: [
          { name: 'Quadriceps', test: 'grading' },
          { name: 'Hamstrings', test: 'grading' },
        ],
      },
      {
        category: 'Special Tests',
        priority: 'high',
        tests: [
          { name: "McMurray's", test: 'binary' },
          { name: "Lachman's", test: 'binary' },
          { name: 'Anterior Drawer', test: 'binary' },
          { name: 'Posterior Drawer', test: 'binary' },
          { name: 'Patellar grind', test: 'binary' },
        ],
      },
      {
        category: 'Palpation',
        priority: 'medium',
        tests: [
          { name: 'Joint line tenderness', test: 'binary' },
          { name: 'Patellar tenderness', test: 'binary' },
          { name: 'Surrounding musculature', test: 'notes' },
        ],
      },
      {
        category: 'Neurological Screen',
        priority: 'medium',
        tests: [
          { name: 'Dermatomes', test: 'notes' },
          { name: 'Myotomes', test: 'grading' },
          { name: 'Reflexes', test: 'notes' },
        ],
      },
    ],
  },
  {
    name: 'Shoulder Pain Assessment',
    complaint: 'Shoulder Pain',
    categories: [
      {
        category: 'Range of Motion',
        priority: 'high',
        tests: [
          { name: 'Flexion', test: 'measurement', unit: 'degrees' },
          { name: 'Abduction', test: 'measurement', unit: 'degrees' },
          { name: 'Internal rotation', test: 'measurement', unit: 'degrees' },
          { name: 'External rotation', test: 'measurement', unit: 'degrees' },
        ],
      },
      {
        category: 'Strength Testing',
        priority: 'high',
        tests: [
          { name: 'Rotator cuff', test: 'grading' },
          { name: 'Deltoid', test: 'grading' },
          { name: 'Biceps', test: 'grading' },
        ],
      },
      {
        category: 'Special Tests',
        priority: 'high',
        tests: [
          { name: 'Hawkins-Kennedy', test: 'binary' },
          { name: "Neer's", test: 'binary' },
          { name: 'Empty Can', test: 'binary' },
          { name: "Speed's", test: 'binary' },
        ],
      },
    ],
  },
  {
    name: 'Lower Back Pain Assessment',
    complaint: 'Lower Back Pain',
    categories: [
      {
        category: 'Posture & Gait',
        priority: 'high',
        tests: [
          { name: 'Postural assessment', test: 'notes' },
          { name: 'Gait observation', test: 'notes' },
        ],
      },
      {
        category: 'Range of Motion',
        priority: 'high',
        tests: [
          { name: 'Flexion', test: 'measurement', unit: 'degrees' },
          { name: 'Extension', test: 'measurement', unit: 'degrees' },
          { name: 'Lateral flexion', test: 'measurement', unit: 'degrees' },
          { name: 'Rotation', test: 'measurement', unit: 'degrees' },
        ],
      },
      {
        category: 'Strength Testing',
        priority: 'high',
        tests: [
          { name: 'Core stability', test: 'grading' },
          { name: 'Hip flexors', test: 'grading' },
          { name: 'Gluteals', test: 'grading' },
        ],
      },
      {
        category: 'Special Tests',
        priority: 'high',
        tests: [
          { name: 'Straight Leg Raise', test: 'binary' },
          { name: 'FABER', test: 'binary' },
          { name: 'FADIR', test: 'binary' },
          { name: 'Slump test', test: 'binary' },
        ],
      },
      {
        category: 'Neurological Screen',
        priority: 'high',
        tests: [
          { name: 'Dermatomes L1-S1', test: 'notes' },
          { name: 'Myotomes', test: 'grading' },
          { name: 'Reflexes', test: 'notes' },
        ],
      },
    ],
  },
  {
    name: 'Neck Pain Assessment',
    complaint: 'Neck Pain',
    categories: [
      {
        category: 'Range of Motion',
        priority: 'high',
        tests: [
          { name: 'Flexion', test: 'measurement', unit: 'degrees' },
          { name: 'Extension', test: 'measurement', unit: 'degrees' },
          { name: 'Rotation left', test: 'measurement', unit: 'degrees' },
          { name: 'Rotation right', test: 'measurement', unit: 'degrees' },
          { name: 'Lateral flexion', test: 'measurement', unit: 'degrees' },
        ],
      },
      {
        category: 'Strength Testing',
        priority: 'high',
        tests: [
          { name: 'Deep neck flexors', test: 'grading' },
          { name: 'Neck extensors', test: 'grading' },
        ],
      },
      {
        category: 'Special Tests',
        priority: 'high',
        tests: [
          { name: "Spurling's", test: 'binary' },
          { name: 'Distraction test', test: 'binary' },
          { name: 'Upper limb tension test', test: 'binary' },
        ],
      },
      {
        category: 'Neurological Screen',
        priority: 'high',
        tests: [
          { name: 'Dermatomes C4-T1', test: 'notes' },
          { name: 'Myotomes', test: 'grading' },
          { name: 'Reflexes', test: 'notes' },
        ],
      },
    ],
  },
]

export const getObjectiveTemplateByComplaint = (complaint: string): ObjectiveTemplate | undefined => {
  return defaultObjectiveTemplates.find(t => t.complaint === complaint)
}
