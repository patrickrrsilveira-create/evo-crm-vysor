import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { GitFork, Send, Timer, Webhook } from 'lucide-react';
import { NodeConfigModal } from './NodeConfigModal';

const meta: Meta<typeof NodeConfigModal> = {
  title: 'Flow Builder/NodeConfigModal',
  component: NodeConfigModal,
  parameters: {
    docs: {
      description: {
        component:
          'Shared modal chrome for any Flow Builder node configuration form. ' +
          'Discriminated-union API: `variant="simple"` (80% of nodes), `variant="tabs"` ' +
          '(basic/advanced split), or `variant="disclosure"` (body + collapsible advanced ' +
          'settings). Focus trap + ESC + ARIA come from the underlying `@evoapi/design-system` ' +
          'Dialog primitive (Radix). Lifted state mandatory — every prop is controlled by the ' +
          'consumer; the component itself holds no useState.\n\n' +
          'The body / tab / advanced slots below are intentionally abstract placeholders — ' +
          'real consumers (e.g. EVO-1271 for Trigger Event) pass their own form JSX. This ' +
          'card only ships the chrome.',
      },
    },
  },
  argTypes: {
    variant: {
      control: 'select',
      options: ['simple', 'tabs', 'disclosure'],
    },
    dirty: { control: 'boolean' },
    loading: { control: 'boolean' },
    contentClassName: {
      control: 'text',
      description: 'Forwarded onto Dialog.Content className. Use for width overrides.',
    },
  },
};

export default meta;

type Story = StoryObj<typeof NodeConfigModal>;

function PlaceholderBody({ label = 'body slot' }: { label?: string }) {
  return (
    <div className="rounded-md border border-dashed border-flow-panel-divider bg-flow-palette-surface p-6 text-center text-sm text-muted-foreground">
      <code className="text-xs">[Consumer provides this {label} via props]</code>
    </div>
  );
}

export const SimpleDirty: Story = {
  name: 'Simple — dirty (Save enabled)',
  render: function Render() {
    const [open, setOpen] = useState(true);
    return (
      <NodeConfigModal
        open={open}
        onCancel={() => setOpen(false)}
        onSave={() => setOpen(false)}
        variant="simple"
        title="Send message"
        description="Configure the message that will be sent on this branch."
        icon={<Send className="h-5 w-5 text-flow-node-action-message-fg" />}
        dirty
      >
        <PlaceholderBody />
      </NodeConfigModal>
    );
  },
};

export const SimplePristine: Story = {
  name: 'Simple — pristine (Save disabled)',
  render: function Render() {
    const [open, setOpen] = useState(true);
    return (
      <NodeConfigModal
        open={open}
        onCancel={() => setOpen(false)}
        onSave={() => setOpen(false)}
        variant="simple"
        title="Send webhook"
        description="Save stays disabled until the form is dirty."
        icon={<Webhook className="h-5 w-5 text-flow-node-action-webhook-fg" />}
        dirty={false}
      >
        <PlaceholderBody />
      </NodeConfigModal>
    );
  },
};

export const SimpleLoading: Story = {
  name: 'Simple — loading (spinner + disabled, "Saving…" for SR)',
  render: function Render() {
    return (
      <NodeConfigModal
        open
        onCancel={() => undefined}
        onSave={() => undefined}
        variant="simple"
        title="Send message"
        description="Both Save and Cancel are disabled while loading. Screen readers hear “Saving…”."
        icon={<Send className="h-5 w-5 text-flow-node-action-message-fg" />}
        dirty
        loading
      >
        <PlaceholderBody />
      </NodeConfigModal>
    );
  },
};

export const Tabs: Story = {
  name: 'Tabs — basic / advanced',
  render: function Render() {
    const [open, setOpen] = useState(true);
    return (
      <NodeConfigModal
        open={open}
        onCancel={() => setOpen(false)}
        onSave={() => setOpen(false)}
        variant="tabs"
        title="Trigger event"
        description="Basic + advanced split via Radix Tabs. Tab state is uncontrolled here; pass `value` / `onTabChange` to control it from outside."
        icon={<GitFork className="h-5 w-5 text-flow-node-trigger-fg" />}
        dirty
        tabs={[
          { value: 'basic', label: 'Basic', content: <PlaceholderBody label="basic tab body" /> },
          {
            value: 'advanced',
            label: 'Advanced',
            content: <PlaceholderBody label="advanced tab body" />,
          },
        ]}
      />
    );
  },
};

export const Disclosure: Story = {
  name: 'Disclosure — body + collapsible advanced',
  render: function Render() {
    const [open, setOpen] = useState(true);
    return (
      <NodeConfigModal
        open={open}
        onCancel={() => setOpen(false)}
        onSave={() => setOpen(false)}
        variant="disclosure"
        title="Wait condition"
        description="Default fields visible; advanced exposes optional filters."
        icon={<Timer className="h-5 w-5 text-flow-node-control-fg" />}
        dirty
        advanced={<PlaceholderBody label="advanced disclosure body" />}
      >
        <PlaceholderBody label="primary body" />
      </NodeConfigModal>
    );
  },
};

export const DisclosureOpen: Story = {
  name: 'Disclosure — advanced open by default',
  render: function Render() {
    return (
      <NodeConfigModal
        open
        onCancel={() => undefined}
        onSave={() => undefined}
        variant="disclosure"
        title="Wait condition"
        icon={<Timer className="h-5 w-5 text-flow-node-control-fg" />}
        defaultAdvancedOpen
        advancedLabel="Show fewer fields"
        dirty
        advanced={<PlaceholderBody label="advanced disclosure body" />}
      >
        <PlaceholderBody label="primary body" />
      </NodeConfigModal>
    );
  },
};

export const WideContent: Story = {
  name: 'Simple — custom width via contentClassName',
  render: function Render() {
    const [open, setOpen] = useState(true);
    return (
      <NodeConfigModal
        open={open}
        onCancel={() => setOpen(false)}
        onSave={() => setOpen(false)}
        variant="simple"
        title="Send webhook"
        description="contentClassName=max-w-4xl widens the dialog beyond the default max-w-2xl."
        icon={<Webhook className="h-5 w-5 text-flow-node-action-webhook-fg" />}
        contentClassName="max-w-4xl"
        dirty
      >
        <PlaceholderBody />
      </NodeConfigModal>
    );
  },
};
