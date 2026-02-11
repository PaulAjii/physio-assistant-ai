import { Injectable } from '@nestjs/common';
import { createPhysioAgent } from '../mastra/agents/physio-agent';
import { MedicalNoteSchema } from 'src/mastra/schemas/medical-schema';
import * as path from 'path';

@Injectable()
export class MedicalAgentService {
  async processConsultation(
    audioBuffer: Buffer,
    incomingMimeType: string,
    fileName: string,
  ) {
    const agent = createPhysioAgent();

    let mimeType = incomingMimeType;

    if (mimeType === 'application/octet-stream' || !mimeType) {
      const ext = path.extname(fileName).toLowerCase();

      const mimeMap: Record<string, string> = {
        '.mp3': 'audio/mp3',
        '.wav': 'audio/wav',
        '.aac': 'audio/aac',
        '.m4a': 'audio/mp4',
        '.flac': 'audio/flac',
        '.ogg': 'audio/ogg',
        '.webm': 'audio/webm',
      };

      if (mimeMap[ext]) {
        console.log(
          `⚠️ Detected generic MIME. Inferred '${mimeMap[ext]}' from extension '${ext}'`,
        );
        mimeType = mimeMap[ext];
      } else {
        console.warn(
          `⚠️ Could not infer type from '${ext}'. Defaulting to 'audio/aac'`,
        );
        mimeType = 'audio/aac';
      }
    }

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
      const response = await agent.generate([{ role: 'user', content }], {
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
      console.error('Error processing consultation: ', err);
      throw new Error('Failed to process the consultation.');
    }
  }
}
