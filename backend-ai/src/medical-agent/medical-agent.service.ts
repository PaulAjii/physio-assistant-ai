import { Injectable } from '@nestjs/common';
import { PhysioAgent } from '../mastra/agents/physio-agent';
import { Agent } from '@mastra/core';
import { MedicalNoteSchema } from 'src/mastra/schemas/medical-schema';

@Injectable()
export class MedicalAgentService {
  private agent: Agent;

  constructor() {
    this.agent = PhysioAgent;
  }

  async processConsultation(audioBuffer: Buffer, mimeType: string) {
    const content = [
      {
        type: 'text' as const,
        text: `
SYSTEM STYLE DIRECTIVE:
Output must be clinically concise and documentation-ready.
Avoid narrative tone.
Avoid repetition.
`,
      },
      {
        type: 'text' as const,
        text: `
### ROLE
You are an expert Medical Scribe. Extract data into the provided JSON schema.

### DATA EXTRACTION RULES (STRICT):
1. If a field is not mentioned in the audio:
   - For Strings: Return "Not reported"
   - For Arrays: Return []
   - For Numbers: Return null
2. Language: Translate Yoruba/Pidgin to Professional English.
3. Style: Chronological, clinical, bullet-style phrasing.
`,
      },
      {
        type: 'text' as const,
        text: 'Analyze this consultation audio and generate the JSON note.',
      },
      {
        type: 'file' as const,
        data: audioBuffer,
        mimeType: mimeType,
      },
    ];

    try {
      const response = await this.agent.generate([{ role: 'user', content }], {
        structuredOutput: {
          schema: MedicalNoteSchema,
        },
        modelSettings: {
          temperature: 0.3,
          topP: 0.8,
          maxOutputTokens: 8192,
        },
      });
      return response.object;
    } catch (err) {
      console.error('Error processing consultation: ', err.message || err);
      throw new Error('Failed to process the consultation.');
    }
  }
}
