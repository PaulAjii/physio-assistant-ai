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
            duration: z.string().describe("Number of days or months or years since pain has started"),
            location: z.array(z.string())
        }),
        red_flags: z.array(z.string()).describe("Any serious symptoms mentioned (e.g., weight loss, night pain)"),
        associated_symptoms: z.array(z.string()).describe("Other symptoms mentioned (e.g., numbness, weakness)"),
        pmhx: z.array(z.string()).describe("List relevant medical history (e.g., ypertension, peptic ulcer disease, diabetes mellitus"),
        pdhx: z.array(z.string()).describe("List of medications taken"),
        pshx: z.array(z.string()).describe("List of surgeries that have been carried out"),
        fhx: z.string().describe("Family history as described by the patient e.g., type of house, number of wives and children"),
        shx: z.array(z.string()).describe("List of social activities and leisure activities e.g., smoking history, alcohol consumption, hobbies")
    }),
    suggested_objective_focus: z.array(z.string()).describe("List of tests/movements the clinician should check based on the history"),
})