# frozen_string_literal: true

require 'rails_helper'
require 'sidekiq/testing'

# B1: regression guard for Labels::UpdateService.
#
# Pre-fix, the rename did `contact.label_list.remove(old) + .add(new) + save!`
# which mutates the TagList in place without dirty-tracking `label_list`, so
# `Contact#publish_label_changes` returned early — silently breaking AC3/AC4
# for the rename path.
#
# Post-fix the service routes through `update!(label_list: ...)`, so the
# commit hook diffs the change and fires `:contact_label_added/removed`.
RSpec.describe Labels::UpdateService do
  # The EvoFlow listener is globally subscribed and would otherwise enqueue
  # against a real Sidekiq client (and try to reach Redis on CI) when these
  # specs emit `:contact_label_*`. `fake!` keeps the test self-contained.
  around { |ex| Sidekiq::Testing.fake! { ex.run } }

  let(:contact) { Contact.create!(name: 'Hue', email: "hue-#{SecureRandom.hex(4)}@test.com") }

  it 'emits :contact_label_removed for the old name and :contact_label_added for the new name' do
    contact.update!(label_list: ['old-name'])

    # The service iterates `tagged_contacts.find_in_batches`, so the in-memory
    # contact it mutates is NOT the `contact` let. Subscribe globally instead.
    added = []
    removed = []
    listener = Class.new do
      define_method(:contact_label_added)   { |data| added   << data[:data] }
      define_method(:contact_label_removed) { |data| removed << data[:data] }
    end.new
    Wisper.subscribe(listener) do
      described_class.new(new_label_title: 'new-name', old_label_title: 'old-name').perform
    end

    expect(added.map { |d| d[:label_name] }).to include('new-name')
    expect(removed.map { |d| d[:label_name] }).to include('old-name')
    expect(contact.reload.label_list).to contain_exactly('new-name')
  end

  # L-2: rename of a contact that already carries both
  # old and new titles should drop the old one without emitting a spurious
  # add for the already-present new title.
  it 'does not emit :contact_label_added when the new name is already in the list' do
    contact.update!(label_list: %w[old-name new-name])

    added = []
    removed = []
    listener = Class.new do
      define_method(:contact_label_added)   { |data| added   << data[:data] }
      define_method(:contact_label_removed) { |data| removed << data[:data] }
    end.new
    Wisper.subscribe(listener) do
      described_class.new(new_label_title: 'new-name', old_label_title: 'old-name').perform
    end

    expect(removed.map { |d| d[:label_name] }).to include('old-name')
    expect(added.map { |d| d[:label_name] }).not_to include('new-name')
    expect(contact.reload.label_list).to contain_exactly('new-name')
  end

  # M-3: a no-op rename (`old == new`) must short-circuit
  # the service so neither the remove nor the add is published.
  it 'is a no-op when old and new titles are identical' do
    contact.update!(label_list: ['same'])

    added = []
    removed = []
    listener = Class.new do
      define_method(:contact_label_added)   { |data| added   << data[:data] }
      define_method(:contact_label_removed) { |data| removed << data[:data] }
    end.new
    Wisper.subscribe(listener) do
      described_class.new(new_label_title: 'same', old_label_title: 'same').perform
    end

    expect(added).to be_empty
    expect(removed).to be_empty
    expect(contact.reload.label_list).to contain_exactly('same')
  end
end
