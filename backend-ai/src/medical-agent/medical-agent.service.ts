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
        const base64Audio = audioBuffer.toString('base64');

        const content = [
            {
                type: 'text' as const,
                text: 'Analyze this consultation audio and generate the JSON note.',
            },
            {
                type: 'file' as const,
                data: base64Audio, 
                mimeType: mimeType, 
            },
        ];

        try {
            const response = await this.agent.generate([
                { role: 'user', content }
            ],
                { 
                    structuredOutput: {
                        schema: MedicalNoteSchema,
                    }
                }
            )
            return response.object
        } catch(err) {
            console.error('Error processing consultation: ', err.message || err)
            throw new Error('Failed to process the consultation.')
        }
    }
}
