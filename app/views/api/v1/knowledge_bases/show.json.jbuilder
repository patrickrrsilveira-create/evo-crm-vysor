json.id @knowledge_base.id
json.name @knowledge_base.name
json.description @knowledge_base.description
json.created_at @knowledge_base.created_at
json.updated_at @knowledge_base.updated_at
json.documents_count @knowledge_base.knowledge_documents.count
json.agent_bots @knowledge_base.agent_bots do |agent|
  json.id agent.id
  json.name agent.name
end
