import { EventEmitter } from "events";

class MessageBus extends EventEmitter { }

export default new MessageBus();
