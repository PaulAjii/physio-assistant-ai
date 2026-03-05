import type { ObjectiveTemplate } from "~~/types/assessment";

export const defaultObjectiveTemplates: ObjectiveTemplate[] = [
  {
    name: "Knee Pain Assessment",
    complaint: "Knee Pain",
    categories: [
      {
        category: "Gait Analysis",
        priority: "high",
        tests: [
          {
            name: "Observation of gait pattern",
            test: "notes",
          },
          {
            name: "Antalgic gait check",
            test: "binary",
          },
        ],
      },
      {
        category: "Range of Motion Measurement",
        priority: "high",
        tests: [
          {
            name: "Passive Knee Flexion",
            test: "measurement",
            unit: "degrees",
          },
          {
            name: "Passive Knee Flexion range, pain assessment, and end-feel evaluation",
            test: "notes",
          },
          {
            name: "Active Knee Flexion",
            test: "measurement",
            unit: "degrees",
          },
          {
            name: "Active Knee Flexion range, pain assessment, and end-feel evaluation",
            test: "notes",
          },
          {
            name: "Passive Knee Extension",
            test: "measurement",
            unit: "degrees",
          },
          {
            name: "Passive Knee Extension range, pain assessment, and end-feel evaluation",
            test: "notes",
          },
          {
            name: "Active Knee Extension",
            test: "measurement",
            unit: "degrees",
          },
          {
            name: "Active Knee Extension range, pain assessment, and end-feel evaluation",
            test: "notes",
          },
        ],
      },
      {
        category: "Strength Testing",
        priority: "high",
        tests: [
          {
            name: "Quadriceps Strength (Oxford Scale)",
            test: "grading",
          },
          {
            name: "Hamstring Strength (Oxford Scale)",
            test: "grading",
          },
        ],
      },
      {
        category: "Special Tests",
        priority: "medium",
        tests: [
          {
            name: "Anterior Drawer Test",
            test: "binary",
          },
          {
            name: "Lachman's Test",
            test: "binary",
          },
          {
            name: "McMurray's Test",
            test: "binary",
          },
          {
            name: "Patellar Grind Test",
            test: "binary",
          },
          {
            name: "Posterior Drawer Test",
            test: "binary",
          },
        ],
      },
    ],
  }
]

export const getObjectiveTemplateByComplaint = (complaint: string): ObjectiveTemplate | undefined => {
  return defaultObjectiveTemplates.find(t => t.complaint.toLowerCase() === complaint.toLowerCase())
}
