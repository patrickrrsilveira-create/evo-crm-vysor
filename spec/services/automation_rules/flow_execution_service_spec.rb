# frozen_string_literal: true

require 'rails_helper'

# B1: regression guard for AutomationRules::FlowExecutionService.
#
# Pre-fix, `add_label`/`remove_label` on the @contact branch did:
#
#     @contact.label_list.add(titles); @contact.save!
#
# That mutates the cached TagList in place without dirty-tracking the
# `label_list` attribute, so `Contact#publish_label_changes` returned early
# and no `:contact_label_added/removed` Wisper event was fired — silently
# breaking AC3 (path 3) and AC4 for the dominant automation path.
#
# Post-fix, both methods route through `update!(label_list: ...)` so the
# setter dirty-tracks and the commit hook diffs the change. These specs
# subscribe to the contact directly and assert the events fire.
RSpec.describe AutomationRules::FlowExecutionService do
  let(:user) { User.create!(name: 'Agent', email: "agent-#{SecureRandom.hex(4)}@test.com") }
  let(:rule) do
    r = AutomationRule.new(
      name: "rule-#{SecureRandom.hex(4)}",
      event_name: 'contact_updated',
      active: true,
      mode: 'flow',
      conditions: [],
      actions: []
    )
    r.save!(validate: false)
    r
  end
  let(:contact) { Contact.create!(name: 'Lead', email: "lead-#{SecureRandom.hex(4)}@test.com") }
  let(:label_vip) { Label.create!(title: 'vip', color: '#fff') }
  let(:label_beta) { Label.create!(title: 'beta', color: '#000') }
  let(:service) { described_class.new(rule, nil, nil, contact) }

  # `FlowExecutionService#initialize` sets `Current.executed_by = rule` as a
  # side-effect. Since these specs drive `add_label`/`remove_label` directly
  # (bypassing `perform`'s `ensure Current.reset`), reset by hand to prevent
  # leakage into other specs in the same run.
  after { Current.reset }

  describe '#add_label on @contact' do
    it 'emits :contact_label_added for each title added (AC3 path 3)' do
      collected = []
      listener = Class.new do
        define_method(:contact_label_added) { |data| collected << data[:data] }
      end.new
      contact.subscribe(listener)

      service.send(:add_label, [label_vip.id, label_beta.id])

      expect(collected.map { |d| d[:label_name] }).to contain_exactly('vip', 'beta')
      expect(contact.reload.label_list).to contain_exactly('vip', 'beta')
    end

    it 'does not double-add existing labels' do
      contact.update!(label_list: ['vip'])

      service.send(:add_label, [label_vip.id, label_beta.id])

      expect(contact.reload.label_list).to contain_exactly('vip', 'beta')
    end
  end

  describe '#remove_label on @contact' do
    before { contact.update!(label_list: %w[vip beta]) }

    it 'emits :contact_label_removed for each title removed (AC4)' do
      collected = []
      listener = Class.new do
        define_method(:contact_label_removed) { |data| collected << data[:data] }
      end.new
      contact.subscribe(listener)

      service.send(:remove_label, [label_vip.id])

      expect(collected.map { |d| d[:label_name] }).to include('vip')
      expect(contact.reload.label_list).to contain_exactly('beta')
    end
  end
end
