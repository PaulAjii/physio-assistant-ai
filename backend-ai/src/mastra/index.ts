import { PhysioAgent } from "./agents/physio-agent";
import { Mastra } from '@mastra/core';

export const mastra = new Mastra({
    agents: {PhysioAgent},
});