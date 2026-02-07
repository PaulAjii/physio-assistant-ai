import { Module } from '@nestjs/common';
import { MedicalAgentService } from './medical-agent.service';
import { MedicalAgentController } from './medical-agent.controller';

@Module({
  controllers: [MedicalAgentController],
  providers: [MedicalAgentService],
})
export class MedicalAgentModule {}
