import React from 'react';
import { JsonViewer } from '@textea/json-viewer';
import { SLA } from '../../models/Asset.model';

function Asset(asset: SLA) {
  return (
    // eslint-disable-next-line react/destructuring-assignment
    <JsonViewer value={asset} rootName={asset.id} defaultInspectDepth={0} />
  );
}

export default Asset;
