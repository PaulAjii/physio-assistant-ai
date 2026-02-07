import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { MedicalAgentModule } from './medical-agent/medical-agent.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    MedicalAgentModule,
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    })
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
