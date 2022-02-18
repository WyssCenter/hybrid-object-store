// vendor
import { createMachine } from 'xstate';


const NamespaceMachine = createMachine({
  initial: 'idle',
  states: {
    idle: {
      on: {
        SUBMIT: { target: 'loading' }
      }
    },
    loading: {
      on: {
        RESET: {target: 'idle'},
        SUCCESS: { target: 'success' },
        ERROR: { target: 'error' }
      }
    },
    refetching: {
      on: {
        SUCCESS: { target: 'success' },
        ERROR: { target: 'error' }
      }
    },
    success: {
      on: {
        REFETCH: {target: 'refetching'}
      }
    },
    error: {
      on: {
        RESET: {target: 'idle'}
      }
    }
  }
});


export default NamespaceMachine;
