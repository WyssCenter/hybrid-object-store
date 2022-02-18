// vendor
import { createMachine } from 'xstate';


const ModalMachine = createMachine({
  initial: 'idle',
  states: {
    idle: {
      on: {
        SUBMIT: { target: 'processing' }
      }
    },
    processing: {
      on: {
        SUCCESS: { target: 'success' },
        ERROR: { target: 'error' }
      }
    },
    success: {
      on: {
        RESET: {target: 'idle'}
      }
    },
    error: {
      on: {
        TRY_AGAIN: { target: 'idle' }
      }
    }
  }
});


export default ModalMachine;
