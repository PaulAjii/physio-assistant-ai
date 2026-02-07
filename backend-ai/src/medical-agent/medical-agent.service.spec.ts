import { Test, TestingModule } from '@nestjs/testing';
import { MedicalAgentService } from './medical-agent.service';

describe('MedicalAgentService', () => {
  let service: MedicalAgentService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [MedicalAgentService],
    }).compile();

    service = module.get<MedicalAgentService>(MedicalAgentService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
