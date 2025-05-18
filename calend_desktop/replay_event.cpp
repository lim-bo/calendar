#include "replay_event.h"
#include "ui_replay_event.h"

replay_event::replay_event(EventData originalEvent, QWidget *parent)
    : QDialog(parent)
    , ui(new Ui::replay_event), originalEvent(originalEvent)
{
    ui->setupUi(this);
}

replay_event::~replay_event()
{
    delete ui;
}

void replay_event::on_pushButton_clicked()
{
    int daysToAdd = ui->spinBox->value();

    repeatedEvent = originalEvent;

    repeatedEvent.start = originalEvent.start.addDays(daysToAdd);
    repeatedEvent.end = originalEvent.end.addDays(daysToAdd);

    repeatedEvent.ID = "";
    this->accept();
}


EventData replay_event::getRepeatedEvent() const
{
    return repeatedEvent;
}
