// vendor
import React, {
  FC,
  useCallback,
  useContext,
  useState,
} from 'react';
// components
import { SectionCard } from 'Components/card/index';
import SyncTable from './table/SyncTable';

interface Props {
  sync: any;
  namespace: string;
}

const SyncSection: FC<Props> = ({sync, namespace}: Props) => {
  return (
    <section className="SyncSection">
      <h4>Namespace Sync Configuration</h4>
      <SectionCard>
        <SyncTable sync={sync} namespace={namespace} />
      </SectionCard>
    </section>
  );
}

export default SyncSection;
