import { z } from 'zod';

export const MedicalNoteSchema = z.object({
    meta: z.object({
        language_detected: z.string().describe("e.g., 'English', 'Yoruba', 'Mixed'"),
        speaker_identification: z.string().describe("Confidence level in separating clinician vs patient"),
    }),
    subjective: z.object({
        hpc: z.string().describe("History of Presenting Complaint. Use patient's verbatim words where possible."),
        pain_profile: z.object({
            intensity: z.number().nullable().describe("0-10 scale"),
            quality: z.string().describe("e.g., Sharp, Dull, 'Ta', 'Ro'"),
            aggravating: z.array(z.string()),
            alleviating: z.array(z.string()),
        }),
        red_flags: z.array(z.string()).describe("Any serious symptoms mentioned (e.g., weight loss, night pain)"),
        associated_symptoms: z.array(z.string()).describe("Other symptoms mentioned (e.g., numbness, weakness)"),
    }),
    suggested_objective_focus: z.array(z.string()).describe("List of tests/movements the clinician should check based on the history"),
})