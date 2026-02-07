import { Controller, Post, UploadedFile, UseInterceptors } from '@nestjs/common';
import { MedicalAgentService } from './medical-agent.service';
import { FileInterceptor } from '@nestjs/platform-express';
import type { Express } from 'express';

@Controller('ai')
export class MedicalAgentController {
  constructor(private readonly medicalAgentService: MedicalAgentService) {}

  @Post('/process-audio')
  @UseInterceptors(FileInterceptor('file'))
  async processAudio(@UploadedFile() file: Express.Multer.File) {
    if (!file) {
      return { error: 'No file uploaded' };
    }

    const result = await this.medicalAgentService.processConsultation(file.buffer, file.mimetype);
    return { message: 'Audio processed successfully', data: result };
  }
}
