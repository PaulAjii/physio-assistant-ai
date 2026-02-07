import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  getHello(): string {
    return 'Hello World!';
  }

  receiveAudio(): string {
    return 'Received Audio for processing';
  }
}
