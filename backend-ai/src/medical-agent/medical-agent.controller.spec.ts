import { Test, TestingModule } from '@nestjs/testing';
import { MedicalAgentController } from './medical-agent.controller';
import { MedicalAgentService } from './medical-agent.service';

describe('MedicalAgentController', () => {
  let controller: MedicalAgentController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [MedicalAgentController],
      providers: [MedicalAgentService],
    }).compile();

    controller = module.get<MedicalAgentController>(MedicalAgentController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
