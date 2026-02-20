import {
  Controller,
  HttpCode,
  HttpStatus,
  ParseFilePipeBuilder,
  Post,
  UploadedFile,
  UseInterceptors,
} from '@nestjs/common';
import { MedicalAgentService } from './medical-agent.service';
import { FileInterceptor } from '@nestjs/platform-express';
import type { Express } from 'express';

@Controller('ai')
export class MedicalAgentController {
  constructor(private readonly medicalAgentService: MedicalAgentService) {}

  @Post('/process-audio')
  @HttpCode(HttpStatus.OK)
  @UseInterceptors(
    FileInterceptor('file', {
      limits: { fileSize: 50 * 1024 * 1024 },
    }),
  )
  async processAudio(
    @UploadedFile(
      new ParseFilePipeBuilder()
        .addMaxSizeValidator({
          maxSize: 50 * 1024 * 1024,
        })
        .build({
          errorHttpStatusCode: HttpStatus.UNPROCESSABLE_ENTITY,
        }),
    )
    file: Express.Multer.File,
  ) {
    if (!file) {
      return { error: 'No file uploaded' };
    }

    const result = await this.medicalAgentService.processConsultation(
      file.buffer,
      file.mimetype,
      file.originalname,
    );
    return { message: 'Audio processed successfully', data: result };
  }
}
