export interface PainProfile {
  intensity: number;
  quality: string;
  aggravating: string[];
  alleviating: string[];
  duration: string;
  location: string[];
}

export interface Subjective {
  presenting_complaint: string;
  history_of_complaint: string;
  pain_profile: PainProfile;
  red_flags: string[];
  associated_symptoms: string[];
  relevant_medical_history: string[];
  drug_history: string[];
  past_surgical_history: string[];
  family_history: string;
  social_history: string[];
}

export interface ObjectiveFinding {
  category: string;
  test: string;
  result: string;
  notes: string;
}

export interface Assessment {
  presenting_complaint: string;
  history_of_complaint: string;
  pain_profile: PainProfile;
  red_flags: string[];
  associated_symptoms: string[];
  relevant_medical_history: string[];
  drug_history: string[];
  past_surgical_history: string[];
  family_history: string;
  social_history: string[];
  objective_findings: ObjectiveFinding[];
}

export interface AIData {
  meta: {
    language_detected: string;
    speaker_identification: string;
  };
  subjective: Subjective;
  suggested_objective_focus: string[]
}

export interface AIDraft {
  message: string;
  data: AIData;
}

export interface CollatedAssessment {
  id: string;
  complaint: string;
  assessment: Assessment;
  ai_draft: AIDraft;
  created_at: string;
}

export interface AssessmentSubmission {
  jobId: string;
  complaint: string;
  corrected_subjective: Subjective;
  objective_findings: ObjectiveFinding[];
}
